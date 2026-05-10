package repository_test

import (
	"testing"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupPrefDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.Prefecture{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "prefectures") })
	return db
}

func TestPrefectureGetAll_ReturnsOrderedByID(t *testing.T) {
	db := setupPrefDB(t)
	repo := repository.NewPrefectureRepository(db)

	prefs := []model.Prefecture{{Name: "北海道"}, {Name: "青森県"}, {Name: "東京都"}}
	if err := db.CreateInBatches(prefs, len(prefs)).Error; err != nil {
		t.Fatal(err)
	}

	result, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 3 {
		t.Fatalf("want 3 prefectures, got %d", len(result))
	}
	// ordered by id: IDs should be ascending
	for i := 1; i < len(result); i++ {
		if result[i].ID < result[i-1].ID {
			t.Errorf("prefectures not ordered by ID: result[%d].ID=%d < result[%d].ID=%d",
				i, result[i].ID, i-1, result[i-1].ID)
		}
	}
}

func TestPrefectureGetAll_EmptyWhenNone(t *testing.T) {
	repo := repository.NewPrefectureRepository(setupPrefDB(t))

	result, err := repo.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 0 {
		t.Errorf("want 0 prefectures, got %d", len(result))
	}
}
