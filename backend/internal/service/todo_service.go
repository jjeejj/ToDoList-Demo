package service

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/jjeejj/todolist/backend/internal/repository"
	todolistv1 "github.com/jjeejj/todolist/backend/proto/todolist/v1"
	"github.com/jjeejj/todolist/backend/proto/todolist/v1/todolistv1connect"
)

// TodoService implements the TodoService RPC methods
type TodoService struct {
	repo *repository.TaskRepository
}

// NewTodoService creates a new TodoService instance
func NewTodoService(repo *repository.TaskRepository) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

// AddTask adds a new task
func (s *TodoService) AddTask(
	ctx context.Context,
	req *connect.Request[todolistv1.AddTaskRequest],
) (*connect.Response[todolistv1.AddTaskResponse], error) {
	if req.Msg.Text == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task text cannot be empty"))
	}

	task := s.repo.AddTask(req.Msg.Text)

	response := &todolistv1.AddTaskResponse{
		Task:    task.ToProtoTask(),
		Success: true,
	}

	return connect.NewResponse(response), nil
}

// GetTasks retrieves all tasks
func (s *TodoService) GetTasks(
	ctx context.Context,
	req *connect.Request[todolistv1.GetTasksRequest],
) (*connect.Response[todolistv1.GetTasksResponse], error) {
	tasks := s.repo.GetTasks()

	protoTasks := make([]*todolistv1.Task, len(tasks))
	for i, task := range tasks {
		protoTasks[i] = task.ToProtoTask()
	}

	response := &todolistv1.GetTasksResponse{
		Tasks: protoTasks,
	}

	return connect.NewResponse(response), nil
}

// DeleteTask removes a task by ID
func (s *TodoService) DeleteTask(
	ctx context.Context,
	req *connect.Request[todolistv1.DeleteTaskRequest],
) (*connect.Response[todolistv1.DeleteTaskResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task ID cannot be empty"))
	}

	success := s.repo.DeleteTask(req.Msg.Id)

	response := &todolistv1.DeleteTaskResponse{
		Success: success,
	}

	return connect.NewResponse(response), nil
}

// UpdateTask updates a task's completion status
func (s *TodoService) UpdateTask(
	ctx context.Context,
	req *connect.Request[todolistv1.UpdateTaskRequest],
) (*connect.Response[todolistv1.UpdateTaskResponse], error) {
	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("task ID cannot be empty"))
	}

	task, success := s.repo.UpdateTask(req.Msg.Id, req.Msg.Completed)
	if !success {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("task not found"))
	}

	response := &todolistv1.UpdateTaskResponse{
		Task:    task.ToProtoTask(),
		Success: true,
	}

	return connect.NewResponse(response), nil
}

// Ensure TodoService implements the interface
var _ todolistv1connect.TodoServiceHandler = (*TodoService)(nil)