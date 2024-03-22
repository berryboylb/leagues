package auth

type SignUpRequest struct {
	FirstName string `json:"first_name" binding:"required,min=3"`
	LastName string `json:"last_name" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}


type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5"`
}

type OtpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Otp string `json:"otp" binding:"required,min=4"`
}


type Otp struct {
	Otp string `json:"otp" binding:"required,min=4"`
	Password string `json:"password" binding:"required,min=5"`
}



type ForgotPasswordRequest struct {
	Email    string `json:"email" binding:"required,email"`
}