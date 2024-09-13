package repository

import (
	"fmt"

	"gorm.io/gorm"
	"app/domain/model"
	"app/domain/repository"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) repository.UserRepository {
	return &userRepository{conn} // ポインタを返す
}

func (r *userRepository) GetEmailAndPass(email string, password string) (*model.User, error) {
	fmt.Printf("%s", "OK?")
	var user model.User
	query := r.Conn.Where("")
	query = query.Where(model.User{Email: email, Password: password})
	err := query.Find(&user).Error
	if err != nil {
		// https://stackoverflow.com/questions/57465968/cannot-use-nil-value-as-a-return-of-type-struct
		// return model.User{}, echo.ErrNotFound
		return nil, err
	}
	return &user, err
}

func (r *userRepository) Create(email string, password string) (*model.User, error) {
	var user model.User
	user.Email = email
	user.Password = password
	// err := r.Conn.Create(user)
	// if err != nil {
	// 	return nil, err
	// }
	return &user, nil
}

func (r *userRepository) Update(u *model.User) (*model.User, error) {
	return nil, nil
}

func (r *userRepository) Delete(id uint) (error) {
	return nil
}