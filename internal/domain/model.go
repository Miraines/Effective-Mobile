package domain

import "time"

// Person – финальная сущность, возвращаемая клиенту.
// @Description  Человек c обогащёнными полями.
// @name         Person
type Person struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Age         *int      `json:"age,omitempty"`
	Gender      *string   `json:"gender,omitempty"`
	CountryID   *string   `json:"country_id,omitempty"`
	Probability *float32  `json:"probability,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
