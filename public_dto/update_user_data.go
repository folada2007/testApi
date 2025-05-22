package public_dto

// NewUserData новые данные
// @Description DTO для обновления данных пользователя
type NewUserData struct {
	Name        string `json:"name"`
	Age         int    `json:"age"`
	Surname     string `json:"surname"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
}
