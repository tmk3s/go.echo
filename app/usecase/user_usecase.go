package usecase

import (
	"fmt"
	"time"

	"app/domain/model"
	"app/domain/repository"
)

// UserParams は PUT /api/user のリクエストボディ。
// フィールドが非公開だと echo の Bind が値を設定できないため、すべて公開する。
// また、送られてこなかった項目を空値で上書きしないよう、すべてポインタで受け取る。
type UserParams struct {
	LastName  *string    `json:"last_name"`
	FirstName *string    `json:"first_name"`
	Email     *string    `json:"email"`
	Gender    *int       `json:"gender"`
	BirthDay  *time.Time `json:"birthday"`
	Working   *bool      `json:"working"`
}

// インターフェースは頭大文字
type UserUseCase interface {
	Get(userId uint) (*model.User, error)
	Update(userId uint, params UserParams) error
}

type userUseCase struct {
	repository.UserRepository
}

func NewUserUseCase(r repository.UserRepository) UserUseCase {
	return &userUseCase{r}
}

func (u *userUseCase) Get(userId uint) (*model.User, error) {
	user, err := u.UserRepository.GetById(userId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUseCase) Update(userId uint, params UserParams) error {
	user, err := u.UserRepository.GetById(userId)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("ユーザーが存在しません")
	}

	// 指定された項目だけを更新する(未指定の項目は既存の値を残す)
	if params.Email != nil && *params.Email != "" {
		user.Email = *params.Email
	}
	user.UserInfo.UserId = user.ID
	if params.LastName != nil {
		user.UserInfo.LastName = *params.LastName
	}
	if params.FirstName != nil {
		user.UserInfo.FirstName = *params.FirstName
	}
	if params.Gender != nil {
		user.UserInfo.Gender = *params.Gender
	}
	if params.BirthDay != nil {
		user.UserInfo.BirthDay = params.BirthDay
	}
	if params.Working != nil {
		user.UserInfo.Working = *params.Working
	}

	if _, err := u.UserRepository.Update(user); err != nil {
		return err
	}
	return nil
}
