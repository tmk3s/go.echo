package repository_test

import (
	"testing"

	"gorm.io/gorm"

	"app/domain/model"
	"app/infrastructure/repository"
)

func setupTodoDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := openTestDB(t)
	if err := db.AutoMigrate(&model.Todo{}); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { truncate(t, db, "todos") })
	return db
}

func seedTodo(t *testing.T, db *gorm.DB, userId uint, title string, completed bool) model.Todo {
	t.Helper()
	todo := model.Todo{UserId: userId, Title: title, Completed: completed}
	if err := db.Create(&todo).Error; err != nil {
		t.Fatal(err)
	}
	return todo
}

// --- Add ---

func TestTodoAdd_AssignsID(t *testing.T) {
	repo := repository.NewTodoRepository(setupTodoDB(t))
	todo, err := repo.Add(&model.Todo{UserId: 1, Title: "新タスク"})
	if err != nil {
		t.Fatal(err)
	}
	if todo.ID == 0 {
		t.Error("expected ID to be assigned after Add")
	}
}

func TestTodoAdd_DefaultsCompleted(t *testing.T) {
	repo := repository.NewTodoRepository(setupTodoDB(t))
	todo, err := repo.Add(model.NewTodo(1, "タスク"))
	if err != nil {
		t.Fatal(err)
	}
	if todo.Completed {
		t.Error("new todo should not be completed by default")
	}
}

// --- GetList ---

func TestTodoGetList_FiltersByUser(t *testing.T) {
	db := setupTodoDB(t)
	repo := repository.NewTodoRepository(db)

	seedTodo(t, db, 1, "ユーザー1のタスク", false)
	seedTodo(t, db, 1, "ユーザー1の別タスク", false)
	seedTodo(t, db, 2, "ユーザー2のタスク", false)

	todos, err := repo.GetList(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 2 {
		t.Errorf("want 2 todos for user 1, got %d", len(todos))
	}
	for _, td := range todos {
		if td.UserId != 1 {
			t.Errorf("got todo with UserId=%d, want 1", td.UserId)
		}
	}
}

func TestTodoGetList_EmptyForUnknownUser(t *testing.T) {
	repo := repository.NewTodoRepository(setupTodoDB(t))

	todos, err := repo.GetList(999)
	if err != nil {
		t.Fatal(err)
	}
	if len(todos) != 0 {
		t.Errorf("want 0 todos, got %d", len(todos))
	}
}

// --- GetById ---

func TestTodoGetById_Found(t *testing.T) {
	db := setupTodoDB(t)
	repo := repository.NewTodoRepository(db)
	created := seedTodo(t, db, 1, "検索タスク", false)

	todo, err := repo.GetById(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if todo.Title != "検索タスク" {
		t.Errorf("want 検索タスク, got %s", todo.Title)
	}
}

func TestTodoGetById_NotFound(t *testing.T) {
	repo := repository.NewTodoRepository(setupTodoDB(t))

	_, err := repo.GetById(9999)
	if err == nil {
		t.Error("want error for non-existent todo, got nil")
	}
}

// --- Update ---

func TestTodoUpdate_SetsCompleted(t *testing.T) {
	db := setupTodoDB(t)
	repo := repository.NewTodoRepository(db)
	created := seedTodo(t, db, 1, "完了テスト", false)

	created.Completed = true
	updated, err := repo.Update(&created)
	if err != nil {
		t.Fatal(err)
	}
	if !updated.Completed {
		t.Error("want Completed=true after Update, got false")
	}

	// verify persisted
	got, _ := repo.GetById(created.ID)
	if !got.Completed {
		t.Error("Completed was not persisted to DB")
	}
}

// --- Delete ---

func TestTodoDelete_SoftDeletes(t *testing.T) {
	db := setupTodoDB(t)
	repo := repository.NewTodoRepository(db)
	created := seedTodo(t, db, 1, "削除タスク", false)

	if err := repo.Delete(&created); err != nil {
		t.Fatal(err)
	}

	// GetById should fail after soft-delete
	_, err := repo.GetById(created.ID)
	if err == nil {
		t.Error("want error after delete, got nil (soft delete not working)")
	}
}
