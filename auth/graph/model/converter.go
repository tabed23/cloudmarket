package model

func ConvertToGraphQLUser(userModel UserModel) *User {
	return &User{
		ID:           userModel.ID,
		FirstName:    userModel.FirstName,
		LastName:     userModel.LastName,
		Email:        userModel.Email,
		Token:        nil,
		RefreshToken: nil,
		Role:         userModel.Role,
		CreatedAt:    &userModel.CreatedAt,
		UpdatedAt:    &userModel.UpdatedAt,
	}
}

func ConvertToGraphQLNewUser(newUserModel NewUserModel) *NewUser {
	return &NewUser{
		FirstName: newUserModel.FirstName,
		LastName:  newUserModel.LastName,
		Email:     newUserModel.Email,
		Password:  newUserModel.Password,
		UpdatedAt: nil,
	}
}
func ConvertToUserModel(user User) *UserModel {
	return &UserModel{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Email:        user.Email,
		Password:     user.Password,
		Token:        *user.Token,
		RefreshToken: *user.RefreshToken,
		Role:         user.Role,
		CreatedAt:    *user.CreatedAt,
		UpdatedAt:    *user.UpdatedAt,
	}
}

func ConvertToNewUserModel(newUser NewUser) *NewUserModel {
	return &NewUserModel{
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Email:     newUser.Email,
		Password:  newUser.Password,
	}
}
