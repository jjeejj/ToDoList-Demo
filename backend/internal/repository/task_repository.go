package repository

import (
	"sync"
	"time"

	"github.com/google/uuid"
	todolistv1 "github.com/jjeejj/todolist/backend/proto/todolist/v1"
)

// Task represents a todo task
type Task struct {
	ID        string
	Text      string
	CreatedAt time.Time
	Completed bool
}

// TaskRepository manages tasks in memory
type TaskRepository struct {
	tasks map[string]*Task
	mutex sync.RWMutex
}

// NewTaskRepository creates a new task repository
func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[string]*Task),
	}
}

// AddTask adds a new task to the repository
func (r *TaskRepository) AddTask(text string) *Task {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task := &Task{
		ID:        uuid.New().String(),
		Text:      text,
		CreatedAt: time.Now(),
		Completed: false, // 新任务默认未完成
	}

	r.tasks[task.ID] = task
	return task
}

// GetTasks returns all tasks from the repository
func (r *TaskRepository) GetTasks() []*Task {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// DeleteTask removes a task from the repository
func (r *TaskRepository) DeleteTask(id string) bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tasks[id]; exists {
		delete(r.tasks, id)
		return true
	}
	return false
}

// UpdateTask updates a task's completion status
func (r *TaskRepository) UpdateTask(id string, completed bool) (*Task, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if task, exists := r.tasks[id]; exists {
		task.Completed = completed
		return task, true
	}
	return nil, false
}

// ToProtoTask converts internal Task to protobuf Task
func (t *Task) ToProtoTask() *todolistv1.Task {
	return &todolistv1.Task{
		Id:        t.ID,
		Text:      t.Text,
		CreatedAt: t.CreatedAt.Unix(),
		Completed: t.Completed,
	}
}