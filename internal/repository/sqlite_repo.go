package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/Donbassenok/go-lab/internal/model"
)

type SQLitePlantRepo struct {
	db *sql.DB
}

func NewSQLitePlantRepo(db *sql.DB) *SQLitePlantRepo {
	return &SQLitePlantRepo{db: db}
}

func (r *SQLitePlantRepo) Create(plant model.Plant) (int, error) {
	query := `INSERT INTO plants (name, species, age) VALUES (?, ?, ?) RETURNING id`
	var id int
	err := r.db.QueryRow(query, plant.Name, plant.Species, plant.Age).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

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

func (r *SQLitePlantRepo) Patch(id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	query := "UPDATE plants SET "
	var args []interface{}
	var setParts []string

	validColumns := map[string]bool{
		"name":    true,
		"species": true,
		"age":     true,
	}

	for key, value := range updates {
		if validColumns[key] {
			setParts = append(setParts, key+" = ?")
			args = append(args, value)
		}
	}

	if len(setParts) == 0 {
		return nil
	}

	query += strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, id)

	res, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("рослину не знайдено для часткового оновлення")
	}
	
	return nil
}
