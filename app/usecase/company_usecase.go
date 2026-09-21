package usecase

import (
	"fmt"
	"strings"

	"app/domain/model"
	"app/domain/repository"
)

// ValidationError は利用者に見せてよい入力不備を表す。
// ハンドラはこれを 400、それ以外(DB障害など)を 500 として扱う。
// 未認証で叩けるエンドポイントで生の DB エラーを返すと、
// テーブル名やインデックス名、登録済みかどうかが漏れるため区別する。
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string { return e.Message }

func invalid(format string, args ...any) error {
	return &ValidationError{Message: fmt.Sprintf(format, args...)}
}

// CompanyRegistrationInput は会社登録(未ログインで実行可能)の入力。
type CompanyRegistrationInput struct {
	CompanyName string
	Email       string
	Password    string
	LastName    string
	FirstName   string
}

const minPasswordLength = 8

type CompanyUseCase interface {
	GetCompany(id uint) (*model.Company, error)
	UpdateCompany(id uint, name string) error
	Register(input CompanyRegistrationInput) (*model.Company, error)
}

type companyUseCase struct {
	repo     repository.CompanyRepository
	userRepo repository.UserRepository
}

func NewCompanyUseCase(r repository.CompanyRepository, userRepo repository.UserRepository) CompanyUseCase {
	return &companyUseCase{repo: r, userRepo: userRepo}
}

func (u *companyUseCase) GetCompany(id uint) (*model.Company, error) {
	return u.repo.GetById(id)
}

func (u *companyUseCase) UpdateCompany(id uint, name string) error {
	company, err := u.repo.GetById(id)
	if err != nil {
		return err
	}
	company.Name = name
	return u.repo.Update(company)
}

// Register は会社と、その会社に所属する最初のユーザーを作成する。
// 作成されたユーザーは company_id を持つため、そのままログインして利用できる。
func (u *companyUseCase) Register(input CompanyRegistrationInput) (*model.Company, error) {
	companyName := strings.TrimSpace(input.CompanyName)
	email := strings.TrimSpace(input.Email)

	if companyName == "" {
		return nil, invalid("会社名を入力してください")
	}
	if email == "" {
		return nil, invalid("メールアドレスを入力してください")
	}
	if len(input.Password) < minPasswordLength {
		return nil, invalid("パスワードは%d文字以上で入力してください", minPasswordLength)
	}
	if strings.TrimSpace(input.LastName) == "" || strings.TrimSpace(input.FirstName) == "" {
		return nil, invalid("姓と名を入力してください")
	}

	existing, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, invalid("このメールアドレスはすでに登録されています")
	}

	company := &model.Company{Name: companyName}
	user := model.NewCompanyOwner(email, input.Password, strings.TrimSpace(input.LastName), strings.TrimSpace(input.FirstName))
	if err := u.repo.Register(company, user); err != nil {
		return nil, err
	}
	return company, nil
}
