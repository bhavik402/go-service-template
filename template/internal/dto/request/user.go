package request

type CreateUser struct {
	Name  string `json:"name" validate:"required,min=1,max=200"`
	Email string `json:"email" validate:"required,email"`
}

type UpdateUser struct {
	Name  string `json:"name" validate:"required,min=1,max=200"`
	Email string `json:"email" validate:"required,email"`
}
