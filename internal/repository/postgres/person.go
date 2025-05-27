package postgres

import (
	"Effective-Mobile/internal/domain"
	"Effective-Mobile/internal/repository"
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type personRepo struct {
	db *pgxpool.Pool
}

func NewPersonRepo(db *pgxpool.Pool) repository.PersonRepository {
	return &personRepo{db: db}
}

func (r *personRepo) Create(ctx context.Context, p *domain.Person) (int64, error) {
	const q = `
		INSERT INTO people (name, surname, patronymic, age, gender, country_id, probability)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at;
	`
	err := r.db.QueryRow(ctx, q,
		p.Name, p.Surname, p.Patronymic,
		p.Age, p.Gender, p.CountryID, p.Probability).
		Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
	return p.ID, err
}

func (r *personRepo) GetByID(ctx context.Context, id int64) (*domain.Person, error) {
	const q = `SELECT id, name, surname, patronymic, age, gender, country_id,
                      probability, created_at, updated_at
               FROM people WHERE id=$1`
	var p domain.Person
	err := r.db.QueryRow(ctx, q, id).
		Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic,
			&p.Age, &p.Gender, &p.CountryID, &p.Probability,
			&p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *personRepo) List(ctx context.Context, f repository.ListFilter) ([]domain.Person, int, error) {
	var (
		args   []any
		where  []string
		people []domain.Person
	)

	if f.Name != "" {
		args = append(args, "%"+f.Name+"%")
		where = append(where, fmt.Sprintf("name ILIKE $%d", len(args)))
	}
	if f.Surname != "" {
		args = append(args, "%"+f.Surname+"%")
		where = append(where, fmt.Sprintf("surname ILIKE $%d", len(args)))
	}
	if f.AgeFrom != nil {
		args = append(args, *f.AgeFrom)
		where = append(where, fmt.Sprintf("age >= $%d", len(args)))
	}
	if f.AgeTo != nil {
		args = append(args, *f.AgeTo)
		where = append(where, fmt.Sprintf("age <= $%d", len(args)))
	}

	query := `SELECT id, name, surname, patronymic, age, gender, country_id,
	             probability, created_at, updated_at
	          FROM people`
	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}
	if f.SortBy != "" {
		query += " ORDER BY " + f.SortBy
	}
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", f.Limit, f.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Person
		if err = rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic,
			&p.Age, &p.Gender, &p.CountryID, &p.Probability,
			&p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		people = append(people, p)
	}

	var total int
	cntQ := "SELECT count(*) FROM people"
	if len(where) > 0 {
		cntQ += " WHERE " + strings.Join(where, " AND ")
	}
	if err = r.db.QueryRow(ctx, cntQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	return people, total, nil
}

func (r *personRepo) Update(ctx context.Context, p *domain.Person) error {
	const q = `
		UPDATE people
		SET name=$1, surname=$2, patronymic=$3, age=$4,
		    gender=$5, country_id=$6, probability=$7,
		    updated_at=now()
		WHERE id=$8;
	`
	_, err := r.db.Exec(ctx, q, p.Name, p.Surname, p.Patronymic,
		p.Age, p.Gender, p.CountryID, p.Probability, p.ID)
	return err
}

func (r *personRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM people WHERE id=$1`, id)
	return err
}
