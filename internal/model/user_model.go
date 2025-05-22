package model

// User представляет пользователя в системе.
type User struct {
	Id          int    `json:"id"`
	UserName    string `json:"username"`
	Surname     string `json:"surname"`
	Age         int    `json:"age"`
	Nationality string `json:"nationality"`
	Gender      string `json:"gender"`
}
