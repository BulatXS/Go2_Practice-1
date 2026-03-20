package service

import (
	"fmt"
	"sync/atomic"
)

type TaskService struct {
	store   *Store
	counter uint64
}

func NewTaskService() *TaskService {
	return &TaskService{
		store: NewStore(),
	}
}

func (s *TaskService) Create(title, description, due string) Task {
	id := fmt.Sprintf("t_%03d", atomic.AddUint64(&s.counter, 1))

	task := Task{
		ID:          id,
		Title:       title,
		Description: description,
		DueDate:     due,
		Done:        false,
	}

	s.store.Create(task)
	return task
}

func (s *TaskService) GetAll() []Task {
	return s.store.GetAll()
}

func (s *TaskService) GetByID(id string) (Task, bool) {
	return s.store.Get(id)
}

func (s *TaskService) Update(id string, title *string, description *string, dueDate *string, done *bool) (Task, bool) {
	task, ok := s.store.Get(id)
	if !ok {
		return Task{}, false
	}

	if title != nil {
		task.Title = *title
	}
	if description != nil {
		task.Description = *description
	}
	if dueDate != nil {
		task.DueDate = *dueDate
	}
	if done != nil {
		task.Done = *done
	}

	return s.store.Update(id, task)
}

func (s *TaskService) Delete(id string) bool {
	return s.store.Delete(id)
}
