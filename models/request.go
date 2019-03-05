package models

// RequestEditUserForm request body for UserRegister
type RequestEditUserForm struct {
	Name     string `json:"name" example:"Alpha Beta" format:"string"`
	Email    string `json:"email" example:"nobody@example.com" format:"email"`
	Password string `json:"password" example:"********" format:"password"`
}

// RequestLoginForm request body for UserLogin
type RequestLoginForm struct {
	Email    string `json:"email" example:"nobody@example.com" format:"email"`
	Password string `json:"password" example:"********" format:"password"`
}
