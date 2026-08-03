// Package handler is the HTTP transport layer: parse requests, call a
// service, write a response. No business logic lives here — if you find
// yourself reaching for an if-statement about business rules, it belongs
// in service/ instead.
package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"{{ module_name }}/internal/apperror"
	"{{ module_name }}/internal/dto/request"
	"{{ module_name }}/internal/dto/response"
	"{{ module_name }}/internal/service"
	"{{ module_name }}/internal/util"
)

type UserHandler struct {
	service  *service.UserService
	validate *validator.Validate
}

func NewUserHandler(svc *service.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{service: svc, validate: validate}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, apperror.InvalidInput("invalid request body"))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		util.WriteError(w, apperror.InvalidInput(err.Error()))
		return
	}

	user, err := h.service.Create(r.Context(), service.CreateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		util.WriteError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusCreated, response.FromDomainUser(user))
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.WriteError(w, apperror.InvalidInput("invalid user id"))
		return
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		util.WriteError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, response.FromDomainUser(user))
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	page := util.ParsePage(r)

	users, total, err := h.service.List(r.Context(), page.Limit(), page.Offset())
	if err != nil {
		util.WriteError(w, err)
		return
	}

	items := make([]response.User, 0, len(users))
	for _, u := range users {
		items = append(items, response.FromDomainUser(u))
	}

	util.WriteJSON(w, http.StatusOK, response.UserList{
		Items: items,
		Meta:  util.NewPageMeta(page, total),
	})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.WriteError(w, apperror.InvalidInput("invalid user id"))
		return
	}

	var req request.UpdateUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.WriteError(w, apperror.InvalidInput("invalid request body"))
		return
	}
	if err := h.validate.Struct(req); err != nil {
		util.WriteError(w, apperror.InvalidInput(err.Error()))
		return
	}

	user, err := h.service.Update(r.Context(), id, service.UpdateUserInput{
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		util.WriteError(w, err)
		return
	}

	util.WriteJSON(w, http.StatusOK, response.FromDomainUser(user))
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		util.WriteError(w, apperror.InvalidInput("invalid user id"))
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		util.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
