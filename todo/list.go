package todo

import "sync"

type List struct {
	tasks map[string]Task
	mtx   sync.RWMutex
}

func NewList() *List {
	return &List{
		tasks: make(map[string]Task),
	}
}

func (l *List) AddTask(task Task) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	if _, ok := l.tasks[task.Title]; ok {
		return ErrTaskAlreadyExist
	}

	l.tasks[task.Title] = task
	return nil
}

func (l *List) ListTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	temp := make(map[string]Task, len(l.tasks))

	for key, value := range l.tasks {
		temp[key] = value
	}
	return temp
}

func (l *List) CompleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}
	task.Complete()

	l.tasks[title] = task

	return nil
}

func (l *List) UncompleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	task, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}
	task.Uncomplete()

	l.tasks[title] = task
	return nil
}

func (l *List) DeleteTask(title string) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	_, ok := l.tasks[title]
	if !ok {
		return ErrTaskNotFound
	}

	delete(l.tasks, title)

	return nil
}

func (l *List) ListNotCompletedTasks() map[string]Task {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	notCompletedTasks := make(map[string]Task)
	for title, task := range l.tasks {
		if !task.Completed {
			notCompletedTasks[title] = task
		}
	}
	return notCompletedTasks
}

func (l *List) GetTask(title string) (Task, error) {
	l.mtx.RLock()
	defer l.mtx.RUnlock()
	task, ok := l.tasks[title]
	if !ok {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}
