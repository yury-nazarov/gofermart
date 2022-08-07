package auth


// User структура для JSON Unmarshal из HTTP Request
type User struct {
	Login string `json:"login,omitempty" binding:"required"`
	Password string `json:"password" binding:"required"`
}
