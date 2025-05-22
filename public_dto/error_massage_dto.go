package public_dto

// Errors представляет собой конкретную ошибку, возникшую при обработке запроса.
// @Description DTO шаблон ошибок
type Errors struct {
	ErrorsType   string `json:"errors_type"`
	ErrorMassage string `json:"error_massage"`
}

// ResponseError представляет собой список ошибок, возвращаемых API.
// @Description DTO для отправки пользователю информацию об ошибках
type ResponseError struct {
	Errors []Errors `json:"errors"`
}
