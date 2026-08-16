package service

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"gorm.io/gorm"

	"github.com/ld/studyroom/internal/constants"
	"github.com/ld/studyroom/internal/model"
	"github.com/ld/studyroom/internal/repository"
	"github.com/ld/studyroom/internal/util"
)

type fakeUserRepo struct{ users map[string]*model.User }

func newFakeUserRepo() *fakeUserRepo { return &fakeUserRepo{users: map[string]*model.User{}} }
func (f *fakeUserRepo) Create(u *model.User) error {
	if _, ok := f.users[u.Username]; ok { return util.ErrConflict }
	f.users[u.Username] = u
	return nil
}
func (f *fakeUserRepo) FindByUsername(username string) (*model.User, error) {
	if u, ok := f.users[username]; ok { return u, nil }
	return nil, util.ErrNotFound
}
func (f *fakeUserRepo) FindByID(id uint) (*model.User, error) {
	for _, u := range f.users { if u.ID == id { return u, nil } }
	return nil, util.ErrNotFound
}
func (f *fakeUserRepo) FindByIDTx(tx *gorm.DB, id uint) (*model.User, error) { return f.FindByID(id) }
func (f *fakeUserRepo) Update(u *model.User) error { f.users[u.Username] = u; return nil }
func (f *fakeUserRepo) UpdateTx(tx *gorm.DB, u *model.User) error { return f.Update(u) }
func (f *fakeUserRepo) ListRanking(period string, limit int) ([]model.User, error) { return nil, nil }

type nilUserRepo struct{}
func (m *nilUserRepo) Create(u *model.User) error { return nil }
func (m *nilUserRepo) FindByUsername(username string) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) FindByID(id uint) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) FindByIDTx(tx *gorm.DB, id uint) (*model.User, error) { return nil, nil }
func (m *nilUserRepo) Update(u *model.User) error { return nil }
func (m *nilUserRepo) UpdateTx(tx *gorm.DB, u *model.User) error { return nil }
func (m *nilUserRepo) ListRanking(period string, limit int) ([]model.User, error) { return nil, nil }

func newUserSvc(repo repository.UserRepository) UserService {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return NewUserService(repo, logger, "test-secret", 72)
}

func TestGetByIDNilReturnsError(t *testing.T) {
	svc := newUserSvc(&nilUserRepo{})
	user, err := svc.GetByID(1)
	if err == nil {
		t.Fatalf("expected error, got user=%+v", user)
	}
}
func TestLoginNilReturnsUnauthorized(t *testing.T) {
	svc := newUserSvc(&nilUserRepo{})
	_, _, err := svc.Login("alice", "123456")
	if err == nil || !errors.Is(err, util.ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}
func TestRegisterDuplicateConflict(t *testing.T) {
	svc := newUserSvc(newFakeUserRepo())
	if _, err := svc.Register("alice", "123456", "测试", "", constants.RoleUser); err != nil {
		t.Fatalf("first register failed: %v", err)
	}
	_, err := svc.Register("alice", "123456", "测试", "", constants.RoleUser)
	if err == nil || !errors.Is(err, util.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
