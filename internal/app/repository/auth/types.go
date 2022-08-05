package auth


// User структура для JSON Unmarshal из HTTP Request
type User struct {
	Login string `json:"login"`
	Password string `json:"password"`
}
