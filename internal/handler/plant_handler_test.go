package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Donbassenok/go-lab/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPlantRepo struct {
	mock.Mock
}

func (m *MockPlantRepo) Create(plant model.Plant) (int, error) {
	args := m.Called(plant)
	return args.Int(0), args.Error(1)
}

func (m *MockPlantRepo) GetAll() ([]model.Plant, error) {
	args := m.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]model.Plant), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockPlantRepo) GetByID(id int) (model.Plant, error) {
	args := m.Called(id)
	return args.Get(0).(model.Plant), args.Error(1)
}

func (m *MockPlantRepo) Update(id int, plant model.Plant) error {
	args := m.Called(id, plant)
	return args.Error(0)
}

func (m *MockPlantRepo) Patch(id int, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockPlantRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestGetPlantByID_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	expectedPlant := model.Plant{ID: 1, Name: "Rose", Species: "Rosa", Age: 2}
	
	mockRepo.On("GetByID", 1).Return(expectedPlant, nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/plants/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responsePlant model.Plant
	err := json.NewDecoder(rr.Body).Decode(&responsePlant)
	assert.NoError(t, err)
	assert.Equal(t, expectedPlant.Name, responsePlant.Name)
	
	mockRepo.AssertExpectations(t)
}

func TestCreatePlant_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	newPlant := model.Plant{Name: "Cactus", Species: "Cactaceae", Age: 5}
	
	mockRepo.On("Create", mock.AnythingOfType("model.Plant")).Return(42, nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(newPlant)
	req := httptest.NewRequest(http.MethodPost, "/plants", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	
	var response map[string]int
	err := json.NewDecoder(rr.Body).Decode(&response)
	
	assert.NoError(t, err)
	assert.Equal(t, 42, response["id"])
	mockRepo.AssertExpectations(t)
}

func TestCreatePlant_ValidationError(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	invalidPlant := model.Plant{Name: "Cactus", Species: "Cactaceae", Age: -5}
	
	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(invalidPlant)
	req := httptest.NewRequest(http.MethodPost, "/plants", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code) 
	mockRepo.AssertNotCalled(t, "Create") 
}

func TestGetAllPlants_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	expectedPlants := []model.Plant{
		{ID: 1, Name: "Rose", Species: "Rosa", Age: 2},
		{ID: 2, Name: "Tulip", Species: "Tulipa", Age: 1},
	}

	mockRepo.On("GetAll").Return(expectedPlants, nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodGet, "/plants", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responsePlants []model.Plant
	err := json.NewDecoder(rr.Body).Decode(&responsePlants)
	
	assert.NoError(t, err)
	assert.Len(t, responsePlants, 2)
	assert.Equal(t, "Rose", responsePlants[0].Name)
	mockRepo.AssertExpectations(t)
}

func TestUpdatePlant_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	updateData := model.Plant{Name: "Big Rose", Species: "Rosa", Age: 3}
	
	mockRepo.On("Update", 1, mock.AnythingOfType("model.Plant")).Return(nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest(http.MethodPut, "/plants/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestPatchPlant_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	mockRepo.On("Patch", 1, mock.Anything).Return(nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	patchData := map[string]interface{}{"age": 10}
	body, _ := json.Marshal(patchData)
	
	req := httptest.NewRequest(http.MethodPatch, "/plants/1", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeletePlant_Success(t *testing.T) {
	mockRepo := new(MockPlantRepo)
	mockRepo.On("Delete", 1).Return(nil)

	h := NewPlantHandler(mockRepo)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/plants/1", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}