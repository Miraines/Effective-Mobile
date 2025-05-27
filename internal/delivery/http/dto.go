package http

// CreatePersonRequest represents payload for POST /people
// @Description  ФИО будущего пользователя.
// @name         CreatePersonRequest
type CreatePersonRequest struct {
	Name       string  `json:"name"       validate:"required,alpharus"`
	Surname    string  `json:"surname"    validate:"required,alpharus"`
	Patronymic *string `json:"patronymic" validate:"omitempty,alpharus"`
}

// UpdatePersonRequest represents payload for PUT /people/{id}
// @Description  Полное обновление полей ФИО.
// @name         UpdatePersonRequest
type UpdatePersonRequest struct {
	Name       string  `json:"name"       validate:"required,alpharus"`
	Surname    string  `json:"surname"    validate:"required,alpharus"`
	Patronymic *string `json:"patronymic" validate:"omitempty,alpharus"`
}
