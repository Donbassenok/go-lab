package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Donbassenok/go-lab/internal/model"
	"github.com/Donbassenok/go-lab/internal/repository"
	"github.com/go-playground/validator/v10"
)

type PlantHandler struct {
	repo     repository.PlantRepo
	validate *validator.Validate
}

func NewPlantHandler(repo repository.PlantRepo) *PlantHandler {
	return &PlantHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

func (h *PlantHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /plants", h.CreatePlant)
	mux.HandleFunc("GET /plants/{id}", h.GetPlantByID)
}

func (h *PlantHandler) CreatePlant(w http.ResponseWriter, r *http.Request) {
	var plant model.Plant
	if err := json.NewDecoder(r.Body).Decode(&plant); err != nil {
		http.Error(w, "Невірний формат JSON", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(plant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.repo.Create(plant)
	if err != nil {
		http.Error(w, "Помилка збереження в базу", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *PlantHandler) GetPlantByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID має бути числом", http.StatusBadRequest)
		return
	}

	plant, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Рослину не знайдено", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(plant)
}