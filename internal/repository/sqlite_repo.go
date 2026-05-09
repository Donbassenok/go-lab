package repository

import (
	"database/sql"
	"errors"

	"github.com/Donbassenok/go-lab/internal/model"
)

type SQLitePlantRepo struct {
	db *sql.DB
}

func NewSQLitePlantRepo(db *sql.DB) *SQLitePlantRepo {
	return &SQLitePlantRepo{db: db}
}

// 1. CREATE - Додавання нової рослини
func (r *SQLitePlantRepo) Create(plant model.Plant) (int, error) {
	query := `INSERT INTO plants (name, species, age) VALUES (?, ?, ?) RETURNING id`
	var id int
	err := r.db.QueryRow(query, plant.Name, plant.Species, plant.Age).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// 2. GET BY ID - Отримання однієї рослини за ID
func (r *SQLitePlantRepo) GetByID(id int) (model.Plant, error) {
	query := `SELECT id, name, species, age FROM plants WHERE id = ?`
	var p model.Plant
	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Species, &p.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return p, errors.New("рослину не знайдено")
		}
		return p, err
	}
	return p, nil
}

// 3. GET ALL - Отримання списку всіх рослин
func (r *SQLitePlantRepo) GetAll() ([]model.Plant, error) {
	query := `SELECT id, name, species, age FROM plants`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []model.Plant
	for rows.Next() {
		var p model.Plant
		if err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Age); err != nil {
			return nil, err
		}
		plants = append(plants, p)
	}
	return plants, nil
}

// 4. UPDATE - Оновлення даних рослини
func (r *SQLitePlantRepo) Update(id int, plant model.Plant) error {
	query := `UPDATE plants SET name = ?, species = ?, age = ? WHERE id = ?`
	res, err := r.db.Exec(query, plant.Name, plant.Species, plant.Age, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("рослину не знайдено для оновлення")
	}
	return nil
}

// 5. DELETE - Видалення рослини
func (r *SQLitePlantRepo) Delete(id int) error {
	query := `DELETE FROM plants WHERE id = ?`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("рослину не знайдено для видалення")
	}
	return nil
}

// 6. GET BY SPECIES - Пошук рослин за видом
func (r *SQLitePlantRepo) GetBySpecies(species string) ([]model.Plant, error) {
	query := `SELECT id, name, species, age FROM plants WHERE species = ?`
	rows, err := r.db.Query(query, species)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []model.Plant
	for rows.Next() {
		var p model.Plant
		if err := rows.Scan(&p.ID, &p.Name, &p.Species, &p.Age); err != nil {
			return nil, err
		}
		plants = append(plants, p)
	}
	return plants, nil
}
