package repository

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"app/domain/model"
	"app/domain/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) repository.UserRepository {
	return &userRepository{conn}
}

func (r *userRepository) GetById(id uint) (*model.User, error) {
	user := &model.User{}
	err := r.Conn.Preload("UserInfo").First(user, id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailAndPass(email string, password string) (*model.User, error) {
	var user model.User
	if err := r.Conn.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, nil
	}
	return &user, nil
}

// hashPassword は平文のパスワードをハッシュ化する。
// 会社登録でも同じ方式を使うため関数に切り出している。
func hashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func (r *userRepository) Create(user *model.User) (*model.User, error) {
	hashed, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashed
	if err := r.Conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Update(user *model.User) (*model.User, error) {
	// 既定の Save は has-one 関連の FK しか更新しないため、UserInfo の列が保存されない。
	// FullSaveAssociations を有効にして関連のカラムも一緒に更新する。
	if err := r.Conn.Session(&gorm.Session{FullSaveAssociations: true}).Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) error {
	return nil
}
