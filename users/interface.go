package users

type UserRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	LastName  string `json:"last_name" binding:"required,min=3"`
	Email     string `json:"email" binding:"required,email"`
}
