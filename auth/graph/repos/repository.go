package repos

import (
	"context"

	"github.com/tabed23/cloudmarket-auth/graph/model"
)

type Repository interface {
	UserCreation(ctx context.Context, input *model.NewUserModel) (*model.UserModel, error)
	UserByEmail(ctx context.Context, email string) (*model.UserModel, error)
	UserByID(ctx context.Context, id string) (*model.UserModel, error)
	UserByRole(ctx context.Context, role string) ([]*model.UserModel, error)
	UserDelete(ctx context.Context, email string) error
	UserUpdate(ctx context.Context, email string, input *model.NewUserModel) (*model.UserModel, error)
}
