package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"

	"github.com/RJHsiao/api-exercise-go/database"
	"github.com/RJHsiao/api-exercise-go/models"
	"github.com/RJHsiao/api-exercise-go/repository"
	"github.com/RJHsiao/api-exercise-go/utilities"
)

// UserRegister godoc
// @Summary Register new user
// @Description Register new user
// @accept application/json
// @produce text/plain
// @param NewUserForm body models.RequestEditUserForm true "new user object"
// @success 200 {null} null "user register successful."
// @failure 400 {null} null "user info not complete."
// @failure 409 {null} null "email is used."
// @Router /register [post]
func UserRegister(w http.ResponseWriter, req *http.Request) {
	var newUserForm models.RequestEditUserForm
	defer req.Body.Close()
	if err := json.NewDecoder(req.Body).Decode(&newUserForm); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if newUserForm.Name == "" || newUserForm.Email == "" || newUserForm.Password == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	isUsed, err := repository.IsEmailRegistered(newUserForm.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if isUsed {
		http.Error(w, "Conflict", http.StatusConflict)
		return
	}

	err = repository.AddUser(newUserForm)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// UserLogin godoc
// @Summary User login
// @Description User login
// @accept application/json
// @produce application/json
// @param LoginForm body models.RequestLoginForm true "user login form"
// @success 200 {object} models.ResponseLogin "user login successful."
// @failure 400 {null} null "login form not complete."
// @failure 404 {null} null "email and/or password incorrect."
// @Router /login [post]
func UserLogin(w http.ResponseWriter, req *http.Request) {
	var userForm models.RequestLoginForm
	defer req.Body.Close()

	if err := json.NewDecoder(req.Body).Decode(&userForm); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if userForm.Email == "" || userForm.Password == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user, err := repository.FindUserByLoginForm(userForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	var res models.ResponseLogin
	res.SessionKey, err = repository.AssignSessionKey(user.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var resBody []byte
	resBody, err = json.Marshal(res)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

// UserLogout godoc
// @Summary User logout
// @Description User logout
// @produce text/plain
// @success 200 {null} null "logout successful."
// @Security SessionKey
// @Router /logout [post]
func UserLogout(w http.ResponseWriter, req *http.Request) {
	sessionKey := req.Header.Get("Session-Key")
	if sessionKey != "" {
		repository.RevokeSessionKey(sessionKey)
	}
	w.WriteHeader(http.StatusOK)
}

// ShowUserInfo godoc
// @Summary Get user info
// @Description Get user info
// @produce application/json
// @success 200 {object} models.ResponseUserInfo "Successful."
// @failure 401 {null} null "session expired or user not login yet."
// @failure 404 {null} null "Session-Key unavailable."
// @Security SessionKey
// @Router /user [get]
func ShowUserInfo(w http.ResponseWriter, req *http.Request) {
	sessionKey := req.Header.Get("Session-Key")
	if sessionKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := repository.FindSessionByKey(sessionKey)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	var user *database.User
	user, err = repository.FindUserByObjectID(session.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var resBody []byte
	resBody, err = json.Marshal(models.ResponseUserInfo{
		Name:       user.Name,
		Email:      user.Email,
		UpdateTime: fmt.Sprint(time.Unix(int64(user.UpdateTime)/1000, 0)),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBody)
}

// EditUserInfo godoc
// @Summary Modify user info
// @Description Modify user info
// @accept application/json
// @produce text/plain
// @param EditUserForm body models.RequestEditUserForm true "edit user object"
// @success 200 {null} null "Edit successful."
// @failure 400 {null} null "body is empty."
// @failure 401 {null} null "session expired or user not login yet."
// @failure 404 {null} null "Session-Key unavailable."
// @failure 409 {null} null "new email is used."
// @Security SessionKey
// @Router /user [patch]
func EditUserInfo(w http.ResponseWriter, req *http.Request) {
	sessionKey := req.Header.Get("Session-Key")
	if sessionKey == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	session, err := repository.FindSessionByKey(sessionKey)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	var user *database.User
	user, err = repository.FindUserByObjectID(session.UserID)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var editUserForm models.RequestEditUserForm
	defer req.Body.Close()
	if err := json.NewDecoder(req.Body).Decode(&editUserForm); err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	isModifiedFieldExist := false
	if editUserForm.Name != "" {
		user.Name = editUserForm.Name
		isModifiedFieldExist = true
	}
	if editUserForm.Password != "" {
		user.Password = utilities.GetSha256SumFromString(editUserForm.Password)
		isModifiedFieldExist = true
	}

	switch {
	case editUserForm.Email != "" && user.Email == editUserForm.Email:
		isModifiedFieldExist = true
		break
	case editUserForm.Email != "" && user.Email != editUserForm.Email:
		isUsed, err := repository.IsEmailRegistered(editUserForm.Email)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if isUsed {
			http.Error(w, "Conflict", http.StatusConflict)
			return
		}
		user.Email = editUserForm.Email
		isModifiedFieldExist = true
		break
	}

	if !isModifiedFieldExist {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	user.UpdateTime = primitive.DateTime(time.Now().Unix() * 1000)

	err = repository.UpdateUser(*user)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
