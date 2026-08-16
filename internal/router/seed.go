package router

import (
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
)

// seedData 初始化种子数据：用户、座位、成就、公告。
func seedData(db *gorm.DB, logger *slog.Logger) error {
	hash := func(pwd string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			logger.Error("seed hash failed", "error", err)
			return ""
		}
		return string(h)
	}

	users := []model.User{
		{Username: "admin", PasswordHash: hash("admin123"), Nickname: "自习室管理员", Role: constants.RoleAdmin},
		{Username: "student1", PasswordHash: hash("stu123456"), Nickname: "学霸小明", Role: constants.RoleUser, Phone: "13800000001", TotalStudyMinutes: 1200},
		{Username: "student2", PasswordHash: hash("stu123456"), Nickname: "考研小红", Role: constants.RoleUser, Phone: "13800000002", TotalStudyMinutes: 860},
		{Username: "student3", PasswordHash: hash("stu123456"), Nickname: "自律小刚", Role: constants.RoleUser, Phone: "13800000003", TotalStudyMinutes: 540},
	}
	if err := db.Create(&users).Error; err != nil {
		return fmt.Errorf("seed users: %w", err)
	}

	seats := []model.Seat{
		{SeatNo: "A-01", Floor: 1, Zone: constants.ZoneSilent, SeatType: "单人桌", Status: constants.SeatIdle, X: 1, Y: 1},
		{SeatNo: "A-02", Floor: 1, Zone: constants.ZoneSilent, SeatType: "单人桌", Status: constants.SeatIdle, X: 2, Y: 1},
		{SeatNo: "A-03", Floor: 1, Zone: constants.ZoneSilent, SeatType: "单人桌", Status: constants.SeatIdle, X: 3, Y: 1},
		{SeatNo: "B-01", Floor: 1, Zone: constants.ZoneDiscussion, SeatType: "双人桌", Status: constants.SeatIdle, X: 1, Y: 2},
		{SeatNo: "B-02", Floor: 1, Zone: constants.ZoneDiscussion, SeatType: "双人桌", Status: constants.SeatIdle, X: 2, Y: 2},
		{SeatNo: "C-01", Floor: 2, Zone: constants.ZoneWindow, SeatType: "靠窗单座", Status: constants.SeatIdle, X: 1, Y: 1},
		{SeatNo: "C-02", Floor: 2, Zone: constants.ZoneWindow, SeatType: "靠窗单座", Status: constants.SeatIdle, X: 2, Y: 1},
	}
	if err := db.Create(&seats).Error; err != nil {
		return fmt.Errorf("seed seats: %w", err)
	}

	achievements := []model.Achievement{
		{Code: "FIRST_BOOK", Name: "初次体验", ConditionType: "study_hours", Threshold: 1, Icon: "🌱", Description: "累计学习 1 小时"},
		{Code: "HOUR_10", Name: "坚持十天", ConditionType: "study_hours", Threshold: 10, Icon: "🔥", Description: "累计学习 10 小时"},
		{Code: "HOUR_50", Name: "自律达人", ConditionType: "study_hours", Threshold: 50, Icon: "🏆", Description: "累计学习 50 小时"},
		{Code: "STREAK_7", Name: "连续七日", ConditionType: "streak_days", Threshold: 7, Icon: "📅", Description: "连续 7 天签到学习"},
	}
	if err := db.Create(&achievements).Error; err != nil {
		return fmt.Errorf("seed achievements: %w", err)
	}

	announcements := []model.Announcement{
		{Title: "自习室新规：静音区保持安静", Content: "请各位同学在静音区内保持安静，手机调至静音，讨论请前往讨论区。", IsActive: true, PublishedAt: time.Now(), CreatedBy: users[0].ID},
		{Title: "暑期开放时间调整", Content: "7 月 1 日起开放时间调整为 07:00-24:00，请合理安排学习时间。", IsActive: true, PublishedAt: time.Now().Add(-time.Hour), CreatedBy: users[0].ID},
		{Title: "成就徽章系统上线", Content: "学习时长达成条件即可获得徽章，快去打卡学习吧！", IsActive: true, PublishedAt: time.Now().Add(-2 * time.Hour), CreatedBy: users[0].ID},
	}
	if err := db.Create(&announcements).Error; err != nil {
		return fmt.Errorf("seed announcements: %w", err)
	}

	records := []model.StudyRecord{
		{UserID: users[1].ID, StudyDate: time.Now().Format("2006-01-02"), DurationMinutes: 180},
		{UserID: users[2].ID, StudyDate: time.Now().Format("2006-01-02"), DurationMinutes: 120},
		{UserID: users[3].ID, StudyDate: time.Now().Format("2006-01-02"), DurationMinutes: 90},
	}
	if err := db.Create(&records).Error; err != nil {
		return fmt.Errorf("seed study records: %w", err)
	}
	logger.Info("seed data created", "users", len(users), "seats", len(seats), "achievements", len(achievements))
	return nil
}
