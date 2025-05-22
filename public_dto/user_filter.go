package public_dto

import "testApi/internal/model"

// UserFilter фильтры и параметры пагинации для получения пользователей
// @Description DTO параметров фильтрации
type UserFilter struct {
	Username    string `json:"username"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
	Age         int    `json:"age"`
	Surname     string `json:"surname"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
}

// PagedUsers список пользователей и пагинация
// @Description DTO передачи данных о пользователе и пагинации
type PagedUsers struct {
	Users      []model.User `json:"users"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"pageSize"`
	TotalPages int          `json:"totalPages"`
}
