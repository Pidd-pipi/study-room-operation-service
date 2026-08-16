package router

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/ld/studyroom/internal/config"
	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/handler"
	"github.com/ld/studyroom/internal/middleware"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/service"
)

// Router 装配依赖并注册路由。
type Router struct {
	cfg    *config.Config
	logger *slog.Logger
	db     *gorm.DB
}

// NewRouter 初始化数据库、依赖并返回 gin 引擎。
func NewRouter(cfg *config.Config, logger *slog.Logger) (*gin.Engine, error) {
	r := &Router{cfg: cfg, logger: logger}
	if err := r.connectDB(); err != nil {
		return nil, err
	}
	if err := r.migrate(); err != nil {
		return nil, err
	}
	if err := r.seed(); err != nil {
		return nil, err
	}
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.ErrorHandler())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     r.cfg.CORSOriginsSlice(),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	engine.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	api := engine.Group("/api")
	api.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.registerV1(api.Group("/v1"))
	return engine, nil
}

func (r *Router) connectDB() error {
	var db *gorm.DB
	var err error
	for i := 0; i < 30; i++ {
		db, err = gorm.Open(mysql.Open(r.cfg.DSN()), &gorm.Config{Logger: gormlogger.Default.LogMode(gormlogger.Warn), DisableForeignKeyConstraintWhenMigrating: true})
		if err == nil {
			sqlDB, dbErr := db.DB()
			if dbErr == nil {
				sqlDB.SetMaxOpenConns(20)
				sqlDB.SetMaxIdleConns(5)
			}
			r.db = db
			r.logger.Info("database connected", "host", r.cfg.DBHost, "db", r.cfg.DBName)
			return nil
		}
		r.logger.Warn("database not ready, retrying", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("connect database: %w", err)
}

func (r *Router) migrate() error {
	if err := r.db.AutoMigrate(
		&model.User{}, &model.Seat{}, &model.Booking{}, &model.StudyRecord{},
		&model.Achievement{}, &model.UserAchievement{}, &model.ViolationRecord{},
		&model.Announcement{}, &model.AuditLog{},
	); err != nil {
		return fmt.Errorf("auto migrate: %w", err)
	}
	r.logger.Info("database migrated")
	return nil
}

func (r *Router) seed() error {
	var count int64
	if err := r.db.Model(&model.User{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count users for seed: %w", err)
	}
	if count > 0 {
		r.logger.Info("seed skipped, users exist")
		return nil
	}
	if err := seedData(r.db, r.logger); err != nil {
		return fmt.Errorf("seed data: %w", err)
	}
	r.logger.Info("seed data inserted")
	return nil
}

func (r *Router) registerV1(v1 *gin.RouterGroup) {
	userRepo := repository.NewUserRepository(r.db)
	seatRepo := repository.NewSeatRepository(r.db)
	bookingRepo := repository.NewBookingRepository(r.db)
	studyRepo := repository.NewStudyRecordRepository(r.db)
	achievementRepo := repository.NewAchievementRepository(r.db)
	violationRepo := repository.NewViolationRepository(r.db)
	announcementRepo := repository.NewAnnouncementRepository(r.db)
	auditRepo := repository.NewAuditLogRepository(r.db)

	userSvc := service.NewUserService(userRepo, r.logger, r.cfg.JWTSecret, r.cfg.TokenTTLHours)
	seatSvc := service.NewSeatService(seatRepo, r.logger)
	achievementSvc := service.NewAchievementService(achievementRepo, r.logger)
	studySvc := service.NewStudyRecordService(studyRepo, r.logger)
	violationSvc := service.NewViolationService(violationRepo, r.logger)
	bookingSvc := service.NewBookingService(bookingRepo, seatRepo, violationRepo, studyRepo, userSvc, r.db, r.logger)
	announcementSvc := service.NewAnnouncementService(announcementRepo, r.logger)

	userHandler := handler.NewUserHandler(userSvc)
	seatHandler := handler.NewSeatHandler(seatSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)
	studyHandler := handler.NewStudyRecordHandler(studySvc, achievementSvc)
	achievementHandler := handler.NewAchievementHandler(achievementSvc)
	violationHandler := handler.NewViolationHandler(violationSvc)
	announcementHandler := handler.NewAnnouncementHandler(announcementSvc)
	auditHandler := handler.NewAuditLogHandler(auditRepo)

	auth := middleware.AuthRequired(r.cfg)
	authLimiter := middleware.RateLimitStrict(r.cfg)
	adminRoles := []constants.UserRole{constants.RoleAdmin}
	userRoles := []constants.UserRole{constants.RoleAdmin, constants.RoleUser}

	registerAuthRoutes(v1, userHandler, authLimiter)
	registerUserRoutes(v1, userHandler, auth)
	registerSeatRoutes(v1, seatHandler, auth, adminRoles)
	registerBookingRoutes(v1, bookingHandler, auth, adminRoles, authLimiter)
	registerStudyRoutes(v1, studyHandler, auth, adminRoles)
	registerAchievementRoutes(v1, achievementHandler, auth, adminRoles)
	registerViolationRoutes(v1, violationHandler, auth, adminRoles)
	registerAnnouncementRoutes(v1, announcementHandler, auth, adminRoles)
	registerAuditLogRoutes(v1, auditHandler, auth, adminRoles)
	_ = userRoles
}
