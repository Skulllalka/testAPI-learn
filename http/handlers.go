package http

import (
	"API/todo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	todoList *todo.List
}

func NewHTTPHandlers(todoList *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todoList: todoList,
	}
}

func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}
	todoTask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todoList.AddTask(todoTask); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskAlreadyExist) {
			http.Error(w, errDTO.ToString(), http.StatusConflict)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}
	b, err := json.MarshalIndent(w, "", "     ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]
	task, err := h.todoList.GetTask(title)
	if err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "     ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write string: ", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todoList.ListTasks()
	b, err := json.MarshalIndent(tasks, "", "     ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write all tasks: ", err)
		return
	}
}

func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todoList.ListNotCompletedTasks()
	b, err := json.MarshalIndent(uncompletedTasks, "", "    ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write response:", err)
		return
	}
}

func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
		return
	}

	title := mux.Vars(r)["title"]
	if completeDTO.Complete {
		if err := h.todoList.CompleteTask(title); err != nil {
			errDTO := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}
			return
		}
	} else {
		if err := h.todoList.UncompleteTask(title); err != nil {
			errDTO := ErrorDTO{
				Message: err.Error(),
				Time:    time.Now(),
			}
			if errors.Is(err, todo.ErrTaskNotFound) {
				http.Error(w, errDTO.ToString(), http.StatusNotFound)
			} else {
				http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
			}
			return
		}
	}
}

func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := mux.Vars(r)["title"]

	if err := h.todoList.DeleteTask(title); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		if errors.Is(err, todo.ErrTaskNotFound) {
			http.Error(w, errDTO.ToString(), http.StatusNotFound)
		} else {
			http.Error(w, errDTO.ToString(), http.StatusInternalServerError)
		}
	}

	return
}
