package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"

	"github.com/handymesh/hyshAuthService/db/mongodb"
	userModel "github.com/handymesh/hyshAuthService/models/user"
)

var (
	r = chi.NewRouter()
)

func init() {
	// Connect to MongoDB
	mongodb.ConnectToMongo()

	// API
	r.Get("/users", List)
	r.Post("/users", Create)
	r.Put("/users/{userId}", Update)
	r.Delete("/users/{userId}", Delete)
}

func TestUser(t *testing.T) {

	ts := httptest.NewServer(r)
	defer ts.Close()

	// New user
	user := &userModel.User{
		Mail:     "test@mail.com",
		Password: "superPass",
	}
	bodyRequest, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", ts.URL+"/users", bytes.NewReader(bodyRequest))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusOK {
		t.Errorf("[Create user] status code: got %v want %v",
			status, http.StatusOK)
	}

	body := w.Body.Bytes()
	json.Unmarshal(body, &user)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("[Create user] status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Get users
	req, _ = http.NewRequest("GET", ts.URL+"/users", nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusOK {
		t.Errorf("[Get users] status code: got %v want %v",
			status, http.StatusOK)
	}

	// Update user
	user.Mail = "update@mail.com"
	bodyRequest, _ = json.Marshal(user)

	userId := user.Id.Hex()
	req, _ = http.NewRequest("PUT", ts.URL+"/users/"+userId, bytes.NewReader(bodyRequest))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusOK {
		t.Errorf("[Update user] status code: got %v want %v",
			status, http.StatusOK)
	}

	// Incorrect Id
	userId = user.Id.Hex() + "1"
	req, _ = http.NewRequest("PUT", ts.URL+"/users/"+userId, bytes.NewReader(bodyRequest))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("[Update user] status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Not Found Id
	userId = "123123123123123123123123"
	req, _ = http.NewRequest("PUT", ts.URL+"/users/"+userId, bytes.NewReader(bodyRequest))

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check the status code is what we expect.
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("[Update user] status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// Delete user
	userId = user.Id.Hex()
	req, _ = http.NewRequest("DELETE", ts.URL+"/users/"+userId, nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("[Delete user] status code: got %v want %v",
			status, http.StatusOK)
	}
}
