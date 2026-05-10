package usecase_test

import (
	"testing"

	"app/domain/model"
	"app/usecase"
)

func newTodoUC(todoRepo *mockTodoRepo, userRepo *mockUserRepo) usecase.TodoUseCase {
	return usecase.NewTodoUseCase(todoRepo, userRepo)
}

// ---- GetTodos ----

func TestGetTodos_ReturnsListForUser(t *testing.T) {
	todoRepo := &mockTodoRepo{
		todos: []model.Todo{
			{Title: "タスク1"},
			{Title: "タスク2"},
		},
	}
	uc := newTodoUC(todoRepo, &mockUserRepo{})

	result, err := uc.GetTodos(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(*result) != 2 {
		t.Errorf("want 2 todos, got %d", len(*result))
	}
}

// ---- AddTodo ----

func TestAddTodo_UserNotFound_ReturnsError(t *testing.T) {
	todoRepo := &mockTodoRepo{}
	userRepo := &mockUserRepo{user: nil} // user not found

	uc := newTodoUC(todoRepo, userRepo)
	err := uc.AddTodo(999, "タスク")
	if err == nil {
		t.Fatal("want error when user not found, got nil")
	}
}

func TestAddTodo_Success_CallsAdd(t *testing.T) {
	user := &model.User{}
	user.ID = 1
	todoRepo := &mockTodoRepo{}
	userRepo := &mockUserRepo{user: user}

	uc := newTodoUC(todoRepo, userRepo)
	if err := uc.AddTodo(1, "新しいタスク"); err != nil {
		t.Fatal(err)
	}
	if !todoRepo.addCalled {
		t.Error("Add was not called on todo repository")
	}
}

// ---- DoneTodo ----

func TestDoneTodo_TodoNotFound_ReturnsError(t *testing.T) {
	todoRepo := &mockTodoRepo{todo: nil, getErr: errNotFound}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	err := uc.DoneTodo(999)
	if err == nil {
		t.Fatal("want error when todo not found, got nil")
	}
}

func TestDoneTodo_SetsCompletedTrue(t *testing.T) {
	todo := &model.Todo{Title: "未完了タスク", Completed: false}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DoneTodo(1); err != nil {
		t.Fatal(err)
	}
	if todoRepo.updatedTodo == nil {
		t.Fatal("Update was not called")
	}
	if !todoRepo.updatedTodo.Completed {
		t.Error("want Completed=true after DoneTodo, got false")
	}
}

// ---- DeleteTodo ----

func TestDeleteTodo_TodoNotFound_ReturnsError(t *testing.T) {
	todoRepo := &mockTodoRepo{todo: nil, getErr: errNotFound}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	err := uc.DeleteTodo(999)
	if err == nil {
		t.Fatal("want error when todo not found, got nil")
	}
}

func TestDeleteTodo_Success_CallsDelete(t *testing.T) {
	todo := &model.Todo{Title: "削除するタスク"}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DeleteTodo(1); err != nil {
		t.Fatal(err)
	}
	if !todoRepo.deleteCalled {
		t.Error("Delete was not called on todo repository")
	}
}
