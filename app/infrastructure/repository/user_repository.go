package repository

import {
	"fmt"

	"gor.io/gorm"
	"app/domain/model"
	"app/domain/repository"
	"app/domain/service"
}

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) repository.UserRepository {
	return &userRepository(conn) // ポインタを返す
}

func (r *userRepository) GetEmailAndPass(email string, password string) (*mode.User, error) {
	fmt.Printf("%s", "OK?")
	var user model.User
	query = r.Conn.Where("")
	query = query.Where(model.User{Email: email, Password: password})
	err != query.Find(&user).Error
	if err != nil {
		// https://stackoverflow.com/questions/57465968/cannot-use-nil-value-as-a-return-of-type-struct
		return model.User{}, echo.ErrNotFound
	}
	return user, err
}

func (r *userRepository) Create(email string, password string) (*mode.User, error) {
	var user model.User
	user.Email = email
	user.password = password
	err != r.Conn.Create(user)

	if err != nil {
		return nil, err
	}
	return user, err
}

func (r *userRepository) Update(u *model.User) (*mode.User, error) {

}

func (r *userRepository) Delete(id uint) (error) {

}