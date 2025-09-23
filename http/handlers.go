package http

import (
	"API/todo"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
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
	}
	if err := taskDTO.ValidateForCreate(); err != nil {
		errDTO := ErrorDTO{
			Message: err.Error(),
			Time:    time.Now(),
		}
		http.Error(w, errDTO.ToString(), http.StatusBadRequest)
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
			return
		}

	}
	b, err := json.MarshalIndent(w, "", "     ")
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(b); err != nil {
		fmt.Println("failed to write http response:", err)
	}
}

func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandlers) HandleCompleteTask(w http.ResponseWriter, r *http.Request) {

}

func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {

}
