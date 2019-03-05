package models

// ResponseLogin Response body for UserLogin
type ResponseLogin struct {
	SessionKey string `json:"session_key" example:"d4877e7e074a50efbba0e71e30f9985a851b7156159c938c5c4c98ae12efb220" format:"string"`
}

// UserInfo Response body for user info
type ResponseUserInfo struct {
	Name       string `json:"name" example:"Alpha Beta" format:"string"`
	Email      string `json:"email" example:"nobody@example.com" format:"email"`
	UpdateTime string `json:"update_at" example:"1970-01-01 00:00:00.000Z" format:"date-time"`
}
