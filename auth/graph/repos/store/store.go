package store

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tabed23/cloudmarket-auth/graph/model"
	"github.com/tabed23/cloudmarket-auth/graph/repos"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

// UserByEmail implements repos.Repository.
func (s *Store) UserByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	var user model.UserModel
	fmt.Println("Fetching user with email:", email)
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	return &user, nil
}

// UserByID implements repos.Repository.
func (s *Store) UserByID(ctx context.Context, id string) (*model.UserModel, error) {
	var user model.UserModel
	if err := s.db.Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to fetch user by id: %w", err)
	}
	return &user, nil
}

// UserByRole implements repos.Repository.
func (s *Store) UserByRole(ctx context.Context, role string) ([]*model.UserModel, error) {
	var users []*model.UserModel
	if err := s.db.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users by role: %w", err)
	}
	return users, nil
}

// UserCreation implements repos.Repository.
func (s *Store) UserCreation(ctx context.Context, input *model.NewUserModel) (*model.UserModel, error) {
	if s == nil {
		return nil, fmt.Errorf("store is nil")
	}

	if s.db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	if input == nil {
		return nil, fmt.Errorf("input parameter is nil")
	}

	usr, err := s.UserByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking existing user: %w", err)
	}
	if usr != nil {
		return nil, fmt.Errorf("user with email %s already exists", input.Email)
	}

	user := model.UserModel{
		ID:        uuid.NewString(),
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		Role:      input.Role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// UserDelete implements repos.Repository.
func (s *Store) UserDelete(ctx context.Context, email string) error {

	if err := s.db.Where("email = ?", email).Delete(&model.UserModel{}).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// UserUpdate implements repos.Repository.
func (s *Store) UserUpdate(ctx context.Context, email string, input *model.NewUserModel) (*model.UserModel, error) {
	var user model.UserModel
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	user.FirstName = input.FirstName
	user.LastName = input.LastName
	user.Password = input.Password
	user.Role = input.Role
	user.UpdatedAt = time.Now()
	if err := s.db.Save(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return &user, nil
}

func NewStore(db *gorm.DB) repos.Repository {
	return &Store{
		db: db,
	}
}
