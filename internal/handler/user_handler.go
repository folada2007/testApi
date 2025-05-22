package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"testApi/internal/model"
	"testApi/internal/repository"
	"testApi/public_dto"
)

// GetFilteredUsersOrDefault godoc
// @Summary Получить отфильтрованных пользователей
// @Description Возвращает список пользователей с фильтрацией и пагинацией
// @Tags users
// @Accept json
// @Produce json
// @Param username query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param nationality query string false "Национальность"
// @Param gender query string false "Гендер"
// @Param age query int false "возраст"
// @Param page query int false "Номер страницы"
// @Param pageSize query int false "размер страницы"
// @Success 200 {object} public_dto.PagedUsers
// @Failure 500 {object} public_dto.ResponseError
// @Router / [get]
func GetFilteredUsersOrDefault(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filter public_dto.UserFilter

		var ageStr = r.URL.Query().Get("age")
		var pageStr = r.URL.Query().Get("page")
		var pageSizeStr = r.URL.Query().Get("pageSize")

		filter.Username = r.URL.Query().Get("username")
		filter.Nationality = r.URL.Query().Get("nationality")
		filter.Gender = r.URL.Query().Get("gender")
		filter.Surname = r.URL.Query().Get("surname")
		filter.Age, _ = strconv.Atoi(ageStr)
		filter.Page, _ = strconv.Atoi(pageStr)
		filter.PageSize, _ = strconv.Atoi(pageSizeStr)

		if filter.Page < 1 {
			filter.Page = 1
		}
		if filter.PageSize < 1 {
			filter.PageSize = 5
		}

		users, err := repo.GetFilteredUsers(filter)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeServerError(err, w, "ошибка получения пользователей из базы данных")
			return
		}
		err = json.NewEncoder(w).Encode(users)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			writeServerError(err, w, "ошибка при кодировании данных пользователя в json")
			return
		}
		w.WriteHeader(http.StatusOK)
		logrus.Info("данные пользователей успешно отправлены")
	}
}

// PostUserNameHandler godoc
// @Summary Создать нового пользователя
// @Description Принимает имя, фамилию и запрашивает определённые данные из внешних API (возраст, национальность, пол).Сохраняет результат в БД
// @Tags users
// @Accept json
// @Produce json
// @Param user body public_dto.UserRegistration true "Данные для регистрации пользователя"
// @Success 201 {string} string "Пользователь успешно создан!"
// @Failure 400 {object} public_dto.ResponseError "Ошибка валидации или получения данных"
// @Router /PostUserName [post]
func PostUserNameHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user public_dto.UserRegistration

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при декодиовании данных пользователя из формата json")
			return
		}

		err = validateUserRegistration(user)
		if err != nil {
			logrus.WithError(err).Warn("ошибка валидации данных пользователя")
			_ = json.NewEncoder(w).Encode(public_dto.ResponseError{[]public_dto.Errors{{
				"ошибка валидации данных", err.Error(),
			}}})
			return
		}

		age, err := getAge(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при обращении к внешнему API для получения возраста")
			return
		}

		nationality, err := getNationality(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при обращении к внешнему API для получения национальности")
			return
		}

		gender, err := getGender(user.Username)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при обращении к внешнему API для получения гендера")
			return
		}

		var fullUserData = model.User{
			UserName:    user.Username,
			Surname:     user.Surname,
			Age:         age,
			Nationality: nationality,
			Gender:      gender,
		}

		err = repo.AddUser(fullUserData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при записи данных пользователя в БД")
			return
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode("Пользователь успешно создан!")

		logrus.Info("пользователь успешно создан")
	}
}

// UpdateUserDataHandler godoc
// @Summary Обновить данные пользователя
// @Description Обновляет данные пользователя по ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "ID пользователя"
// @Param user body public_dto.NewUserData true "Новые данные пользователя"
// @Success 200 {string} string "Данные успешно обновлены!"
// @Failure 400 {object} public_dto.ResponseError "Ошибка обновления или неверные данные"
// @Router /UpdateUser/{id} [put]
func UpdateUserDataHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUserData public_dto.NewUserData

		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			logrus.WithError(err).Warn("неверный формат ID")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(public_dto.ResponseError{[]public_dto.Errors{{
				"серверная ошибка", err.Error(),
			}}})
			return
		}

		err = json.NewDecoder(r.Body).Decode(&newUserData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при декодиовании данных пользователя из формата json")
			return
		}

		err = validationUpdateUserData(newUserData)
		if err != nil {
			logrus.WithError(err).Warn("ошибка при валидации новых данных пользователя")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(public_dto.ResponseError{[]public_dto.Errors{{
				"Неверный формат данных", err.Error(),
			}}})
			return
		}

		err = repo.UpdateUser(id, newUserData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при обновлении данных пользователя в БД")
			return
		}
		_ = json.NewEncoder(w).Encode("Данные успешно обновлены!")
		logrus.Info("данные пользователя успешно обновлены")
	}
}

// DeleteUserDataHandler godoc
// @Summary Удалить пользователя
// @Description Удаляет пользователя по ID
// @Tags users
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {string} string "Пользователь успешно удален!"
// @Failure 400 {object} public_dto.ResponseError "Ошибка удаления пользователя или неверный ID"
// @Router /DeleteUser/{id} [delete]
func DeleteUserDataHandler(repo repository.UserRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr := vars["id"]

		id, err := strconv.Atoi(idStr)
		if err != nil {
			logrus.WithError(err).Warn("неверный формат ID")
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(public_dto.ResponseError{[]public_dto.Errors{{
				"серверная ошибка", err.Error(),
			}}})
			return
		}

		err = repo.DeleteUser(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			writeServerError(err, w, "ошибка при попытке удалить данные пользователя")
			return
		}
		_ = json.NewEncoder(w).Encode("Пользователь успешно удален!")
		logrus.Info("данные пользователя успешно удалены")
	}
}

func getAge(name string) (int, error) {
	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Debug("делаем запрос к agify.io")

	response, err := http.Get(fmt.Sprintf("https://api.agify.io/?name=%s", name))
	if err != nil {
		return 0, errors.New("api request failed")
	}
	defer response.Body.Close()

	type ageResult struct {
		Age int `json:"age"`
	}
	var result ageResult

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return 0, errors.New("decode failed " + err.Error())
	}
	return result.Age, nil
}

func getNationality(name string) (string, error) {
	type country struct {
		CountryId   string  `json:"country_id"`
		Probability float32 `json:"probability"`
	}

	type NationalityResponse struct {
		Count   int       `json:"count"`
		Name    string    `json:"name"`
		Country []country `json:"country"`
	}

	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Debug("делаем запрос к nationalize.io")

	response, err := http.Get(fmt.Sprintf("https://api.nationalize.io/?name=%s", name))
	if err != nil {
		return "", errors.New("api request failed")
	}
	defer response.Body.Close()

	var result NationalityResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", errors.New("decode failed " + err.Error())
	}

	if len(result.Country) == 0 {
		return "unknown", nil
	}

	return result.Country[0].CountryId, nil
}

func getGender(name string) (string, error) {
	logrus.WithFields(logrus.Fields{
		"name": name,
	}).Debug("делаем запрос к genderize.io")

	response, err := http.Get(fmt.Sprintf("https://api.genderize.io/?name=%s", name))
	if err != nil {
		return "", errors.New("api request failed")
	}
	defer response.Body.Close()

	type genderResult struct {
		Gender string `json:"gender"`
	}

	var result genderResult

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return "", errors.New("decode failed " + err.Error())
	}

	return result.Gender, nil
}

func validateUserRegistration(user public_dto.UserRegistration) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Surname == "" {
		return errors.New("surname is required")
	}
	return nil
}

func validationUpdateUserData(newUserData public_dto.NewUserData) error {
	if newUserData.Name == "" {
		return errors.New("name is required")
	}
	if newUserData.Surname == "" {
		return errors.New("surname is required")
	}
	if newUserData.Age <= 0 {
		return errors.New("uncorrected age")
	}
	if newUserData.Gender == "" {
		return errors.New("gender is required")
	}
	if newUserData.Nationality == "" {
		return errors.New("nationality is required")
	}

	return nil
}

func writeServerError(err error, w http.ResponseWriter, logMassage string) {
	logrus.WithError(err).Error(logMassage)
	_ = json.NewEncoder(w).Encode(public_dto.ResponseError{[]public_dto.Errors{{
		"серверная ошибка", "ошибка на стороне сервера",
	}}})
}
