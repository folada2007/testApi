package public_dto

// UserRegistration содержит данные для регистрации пользователя.
// @Description DTO для регистрации нового пользователя
type UserRegistration struct {
	Username string `json:"username"`
	Surname  string `json:"surname"`
}
