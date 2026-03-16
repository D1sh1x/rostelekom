package handler

import (
	"skilltracker/internal/service"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	service   service.ServiceInterface
	validator *validator.Validate
}

func NewHandler(s service.ServiceInterface) *Handler {
	return &Handler{
		service:   s,
		validator: validator.New(),
	}
}

func (h *Handler) validate(i interface{}) error {
	return h.validator.Struct(i)
}
