package session

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	// "go.mongodb.org/mongo-driver/bson"
	"github.com/go-chi/chi"
	chiMiddleware "github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	pb "github.com/handymesh/hyshAuthService/grpc/mail"
	grpcServer "github.com/handymesh/hyshAuthService/grpc/server"
	"github.com/handymesh/hyshAuthService/handlers/user"
	"github.com/handymesh/hyshAuthService/middleware"
	sessionModel "github.com/handymesh/hyshAuthService/models/session"
	userModel "github.com/handymesh/hyshAuthService/models/user"
	"github.com/handymesh/hyshAuthService/utils"
	"github.com/handymesh/hyshAuthService/utils/crypto"
)

var log = logrus.New()

func init() {
	// Logging =================================================================
	// Setup the logger backend using Sirupsen/logrus and configure
	// it to use a custom JSONFormatter. See the logrus docs for how to
	// configure the backend at github.com/Sirupsen/logrus
	log.Formatter = new(logrus.JSONFormatter)
}

// Routes creates a REST router
func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Captcha)
	r.Use(chiMiddleware.AllowContentType("application/json"))

	r.Get("/debug/{token}", Debug)
	r.Post("/", Login)
	r.Post("/new", Registration)
	r.Post("/recovery", Recovery)
	r.Post("/recovery/{token}", RecoveryByToken)
	r.Post("/refresh", Refresh)
	r.Delete("/", Logout)

	return r
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	var user *userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	var passwordUser = user.Password
	var searchUser = userModel.User{Email: user.Email}
	user, err = userModel.FindOne(searchUser)
	if err != nil {
		utils.Error(w, errors.New(`{"mail":"incorrect mail or password"}`), http.StatusBadRequest)
		return
	}

	isErr := crypto.CheckPasswordHash(passwordUser, user.Password)
	if !isErr {
		utils.Error(w, errors.New(`{"mail":"incorrect mail or password"}`), http.StatusBadRequest)
		return
	}

	// Create JWT token
	tokenString, refreshToken, err := CreateJWTToken()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", tokenString)
	w.WriteHeader(http.StatusCreated)

	type UserOutput struct {
		User struct {
			ID     string `json:"Id"`
			Email  string `json:"Email"`
			Gender string `json:"Gender"`
		} `json:"user"`
		Tokens struct {
			Access  string `json:"access"`
			Refresh string `json:"refresh"`
		} `json:"tokens"`
	}

	var userOutput UserOutput

	userOutput.User.ID = user.Id
	userOutput.User.Email = *user.Email
	userOutput.User.Gender = user.Gender
	userOutput.Tokens.Access = tokenString
	userOutput.Tokens.Refresh = refreshToken

	response := utils.ResponseType{
		Data:    userOutput,
		Status:  http.StatusOK,
		Message: "Success",
	}

	output, err := json.Marshal(response)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	w.Write(output)
}

func Registration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	// And now set a new body, which will simulate the same data we read:
	r.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	user.Create(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var Authorization = r.Header.Get("Authorization")
	if Authorization == "" {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"not auth"`), http.StatusBadRequest)
		return
	}

	token, err := sessionModel.VerifyToken(Authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"token invalid"`), http.StatusBadRequest)
		return
	}

	err = sessionModel.Delete(Authorization)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	response := utils.ResponseType{
		Data:    "",
		Status:  http.StatusOK,
		Message: "Success",
	}

	w.WriteHeader(http.StatusOK)
	output, err := json.Marshal(response)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	w.Write(output)
	return
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var TOKEN_REFRESH = r.Header.Get("Authorization")
	if TOKEN_REFRESH == "" {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"not auth"`), http.StatusBadRequest)
		return
	}

	// Chech REFRESH TOKEN
	status, err := sessionModel.CheckRefreshToken(TOKEN_REFRESH)
	if err != nil || status != true {
		w.WriteHeader(http.StatusUnauthorized)
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	// Create JWT token
	tokenString, refreshToken, err := CreateJWTToken()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	w.Header().Set("Authorization", tokenString)
	type Data struct {
		Tokens struct {
			Access  string `json:"access"`
			Refresh string `json:"refresh"`
		} `json:"tokens"`
	}

	var data Data
	data.Tokens.Access = tokenString
	data.Tokens.Refresh = refreshToken

	response := utils.ResponseType{
		Data:    data,
		Status:  http.StatusCreated,
		Message: "Success",
	}

	output, err := json.Marshal(response)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	w.Write(output)
	return
}

func Recovery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	var user *userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	// search user by mail
	searchUser := userModel.User{}
	searchUser.Email = user.Email
	user, err = userModel.FindOne(searchUser)
	if err != nil {
		utils.Error(w, errors.New(`{"mail":"incorrect mail"}`), http.StatusBadRequest)
		return
	}

	// get refresh token
	recoveryLink, err := sessionModel.NewRecoveryLink(user.Id)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	// Send mail
	conn := grpcServer.GetConnClient()
	c := pb.NewMailClient(conn)
	_, err = c.SendMail(context.Background(), &pb.MailRequest{
		Template: "recovery",
		Mail:     *user.Email,
		Url:      "http://localhost:3000/recovery/" + recoveryLink,
	})
	if err != nil {
		fmt.Printf("error %v", err.Error())
		utils.Error(w, errors.New("\"failed to send message\""), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func RecoveryByToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	var user *userModel.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}

	// check correct a password
	if user.Password != user.PasswordRetry {
		utils.Error(w, errors.New(`{"retryPassword":"incorrect new password"}`), http.StatusBadRequest)
		return
	}

	userId, _ := sessionModel.GetValueByKey(user.RecoveryToken)
	if len(userId) == 0 {
		utils.Error(w, errors.New(`"not found"`), http.StatusBadRequest)
		return
	}

	t, err := primitive.ObjectIDFromHex(userId)
	user.Id = t.String()
	user.Password, _ = crypto.HashPassword(user.Password)

	user, err = userModel.Update(user)
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return
	}
	err = sessionModel.Delete(user.RecoveryToken)
	if err != nil {
		utils.Error(w, errors.New("not found"), http.StatusBadRequest)
		return
	}
	// TODO: Send mail (theme: New password)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func Debug(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{}`))
}
