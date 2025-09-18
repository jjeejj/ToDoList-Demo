package service

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/jjeejj/todolist/backend/internal/repository"
	todolistv1 "github.com/jjeejj/todolist/backend/proto/todolist/v1"
)

func TestTodoService_AddTask(t *testing.T) {
	repo := repository.NewTaskRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	t.Run("successful task addition", func(t *testing.T) {
		req := connect.NewRequest(&todolistv1.AddTaskRequest{
			Text: "Test task",
		})

		resp, err := service.AddTask(ctx, req)
		if err != nil {
			t.Fatalf("AddTask failed: %v", err)
		}

		if !resp.Msg.Success {
			t.Error("Expected success to be true")
		}

		if resp.Msg.Task == nil {
			t.Error("Expected task to be returned")
		}

		if resp.Msg.Task.Text != "Test task" {
			t.Errorf("Expected task text to be 'Test task', got '%s'", resp.Msg.Task.Text)
		}

		if resp.Msg.Task.Id == "" {
			t.Error("Expected task ID to be generated")
		}

	})

	t.Run("empty task text should fail", func(t *testing.T) {
		req := connect.NewRequest(&todolistv1.AddTaskRequest{
			Text: "",
		})

		_, err := service.AddTask(ctx, req)
		if err == nil {
			t.Error("Expected error for empty task text")
		}

		connectErr := err.(*connect.Error)
		if connectErr.Code() != connect.CodeInvalidArgument {
			t.Errorf("Expected InvalidArgument error code, got %v", connectErr.Code())
		}
	})
}

func TestTodoService_GetTasks(t *testing.T) {
	repo := repository.NewTaskRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	t.Run("get empty task list", func(t *testing.T) {
		req := connect.NewRequest(&todolistv1.GetTasksRequest{})

		resp, err := service.GetTasks(ctx, req)
		if err != nil {
			t.Fatalf("GetTasks failed: %v", err)
		}

		if len(resp.Msg.Tasks) != 0 {
			t.Errorf("Expected empty task list, got %d tasks", len(resp.Msg.Tasks))
		}
	})

	t.Run("get tasks after adding some", func(t *testing.T) {
		// Add some tasks first
		task1 := repo.AddTask("Task 1")
		task2 := repo.AddTask("Task 2")

		req := connect.NewRequest(&todolistv1.GetTasksRequest{})

		resp, err := service.GetTasks(ctx, req)
		if err != nil {
			t.Fatalf("GetTasks failed: %v", err)
		}

		if len(resp.Msg.Tasks) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(resp.Msg.Tasks))
		}

		// Verify task details
		found1, found2 := false, false
		for _, task := range resp.Msg.Tasks {
			if task.Id == task1.ID && task.Text == "Task 1" {
				found1 = true
			}
			if task.Id == task2.ID && task.Text == "Task 2" {
				found2 = true
			}
		}

		if !found1 {
			t.Error("Task 1 not found in response")
		}
		if !found2 {
			t.Error("Task 2 not found in response")
		}
	})
}

func TestTodoService_DeleteTask(t *testing.T) {
	repo := repository.NewTaskRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	t.Run("successful task deletion", func(t *testing.T) {
		// Add a task first
		task := repo.AddTask("Task to delete")

		req := connect.NewRequest(&todolistv1.DeleteTaskRequest{
			Id: task.ID,
		})

		resp, err := service.DeleteTask(ctx, req)
		if err != nil {
			t.Fatalf("DeleteTask failed: %v", err)
		}

		if !resp.Msg.Success {
			t.Error("Expected success to be true")
		}

		// Verify task is actually deleted
		tasks := repo.GetTasks()
		for _, taskItem := range tasks {
			if taskItem.ID == task.ID {
				t.Error("Task should have been deleted")
			}
		}
	})

	t.Run("delete non-existent task", func(t *testing.T) {
		req := connect.NewRequest(&todolistv1.DeleteTaskRequest{
			Id: "non-existent-id",
		})

		resp, err := service.DeleteTask(ctx, req)
		if err != nil {
			t.Fatalf("DeleteTask failed: %v", err)
		}

		if resp.Msg.Success {
			t.Error("Expected success to be false for non-existent task")
		}
	})

	t.Run("empty task ID should fail", func(t *testing.T) {
		req := connect.NewRequest(&todolistv1.DeleteTaskRequest{
			Id: "",
		})

		_, err := service.DeleteTask(ctx, req)
		if err == nil {
			t.Error("Expected error for empty task ID")
		}

		connectErr := err.(*connect.Error)
		if connectErr.Code() != connect.CodeInvalidArgument {
			t.Errorf("Expected InvalidArgument error code, got %v", connectErr.Code())
		}
	})
}

func TestTodoService_Integration(t *testing.T) {
	repo := repository.NewTaskRepository()
	service := NewTodoService(repo)
	ctx := context.Background()

	t.Run("full workflow integration test", func(t *testing.T) {
		// 1. Start with empty list
		getReq := connect.NewRequest(&todolistv1.GetTasksRequest{})
		getResp, err := service.GetTasks(ctx, getReq)
		if err != nil {
			t.Fatalf("GetTasks failed: %v", err)
		}
		if len(getResp.Msg.Tasks) != 0 {
			t.Error("Expected empty task list initially")
		}

		// 2. Add multiple tasks
		addReq1 := connect.NewRequest(&todolistv1.AddTaskRequest{Text: "Task 1"})
		addResp1, err := service.AddTask(ctx, addReq1)
		if err != nil {
			t.Fatalf("AddTask 1 failed: %v", err)
		}

		addReq2 := connect.NewRequest(&todolistv1.AddTaskRequest{Text: "Task 2"})
		addResp2, err := service.AddTask(ctx, addReq2)
		if err != nil {
			t.Fatalf("AddTask 2 failed: %v", err)
		}

		// 3. Verify both tasks exist
		getResp, err = service.GetTasks(ctx, getReq)
		if err != nil {
			t.Fatalf("GetTasks failed: %v", err)
		}
		if len(getResp.Msg.Tasks) != 2 {
			t.Errorf("Expected 2 tasks, got %d", len(getResp.Msg.Tasks))
		}

		// 4. Delete one task
		delReq := connect.NewRequest(&todolistv1.DeleteTaskRequest{
			Id: addResp1.Msg.Task.Id,
		})
		delResp, err := service.DeleteTask(ctx, delReq)
		if err != nil {
			t.Fatalf("DeleteTask failed: %v", err)
		}
		if !delResp.Msg.Success {
			t.Error("Expected delete to succeed")
		}

		// 5. Verify only one task remains
		getResp, err = service.GetTasks(ctx, getReq)
		if err != nil {
			t.Fatalf("GetTasks failed: %v", err)
		}
		if len(getResp.Msg.Tasks) != 1 {
			t.Errorf("Expected 1 task after deletion, got %d", len(getResp.Msg.Tasks))
		}
		if getResp.Msg.Tasks[0].Id != addResp2.Msg.Task.Id {
			t.Error("Wrong task remained after deletion")
		}
	})
}
