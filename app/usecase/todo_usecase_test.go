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
	err := uc.DoneTodo(1, 999)
	if err == nil {
		t.Fatal("want error when todo not found, got nil")
	}
}

func TestDoneTodo_SetsCompletedTrue(t *testing.T) {
	todo := &model.Todo{UserId: 1, Title: "未完了タスク", Completed: false}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DoneTodo(1, 1); err != nil {
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
	err := uc.DeleteTodo(1, 999)
	if err == nil {
		t.Fatal("want error when todo not found, got nil")
	}
}

func TestDeleteTodo_Success_CallsDelete(t *testing.T) {
	todo := &model.Todo{UserId: 1, Title: "削除するタスク"}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DeleteTodo(1, 1); err != nil {
		t.Fatal(err)
	}
	if !todoRepo.deleteCalled {
		t.Error("Delete was not called on todo repository")
	}
}

// ---- 所有者チェック ----

func TestDoneTodo_OtherUsersTodo_ReturnsError(t *testing.T) {
	todo := &model.Todo{UserId: 2, Title: "他人のタスク"}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DoneTodo(1, 1); err == nil {
		t.Fatal("want error when completing another user's todo, got nil")
	}
	if todoRepo.updatedTodo != nil {
		t.Error("Update must not be called for another user's todo")
	}
}

func TestDeleteTodo_OtherUsersTodo_ReturnsError(t *testing.T) {
	todo := &model.Todo{UserId: 2, Title: "他人のタスク"}
	todo.ID = 1
	todoRepo := &mockTodoRepo{todo: todo}

	uc := newTodoUC(todoRepo, &mockUserRepo{})
	if err := uc.DeleteTodo(1, 1); err == nil {
		t.Fatal("want error when deleting another user's todo, got nil")
	}
	if todoRepo.deleteCalled {
		t.Error("Delete must not be called for another user's todo")
	}
}

func TestTodo_NotFoundAndOthersTodo_ReturnSameError(t *testing.T) {
	// 存在しない Todo
	missing := newTodoUC(&mockTodoRepo{todo: nil, getErr: errNotFound}, &mockUserRepo{})
	missingErr := missing.DoneTodo(1, 999)
	if missingErr == nil {
		t.Fatal("want error when todo does not exist")
	}

	// 他人の Todo
	othersTodo := &model.Todo{UserId: 2, Title: "他人のタスク"}
	othersTodo.ID = 1
	others := newTodoUC(&mockTodoRepo{todo: othersTodo}, &mockUserRepo{})
	othersErr := others.DoneTodo(1, 1)
	if othersErr == nil {
		t.Fatal("want error when the todo belongs to another user")
	}

	// 応答が違うと id の総当たりで他人の Todo の実在を判別できてしまう
	if missingErr.Error() != othersErr.Error() {
		t.Errorf("errors must be indistinguishable: %q vs %q", missingErr.Error(), othersErr.Error())
	}
}
