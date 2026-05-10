package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

func TestPrefectureGetAll_ReturnsAll(t *testing.T) {
	repo := &mockPrefRepo{
		prefectures: []model.Prefecture{
			{Name: "北海道"},
			{Name: "東京都"},
			{Name: "大阪府"},
		},
	}
	uc := usecase.NewPrefectureUseCase(repo)

	result, err := uc.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Errorf("want 3 prefectures, got %d", len(result))
	}
	if result[0].Name != "北海道" {
		t.Errorf("want 北海道, got %s", result[0].Name)
	}
}

func TestPrefectureGetAll_EmptyWhenNone(t *testing.T) {
	repo := &mockPrefRepo{prefectures: []model.Prefecture{}}
	uc := usecase.NewPrefectureUseCase(repo)

	result, err := uc.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("want 0 prefectures, got %d", len(result))
	}
}
