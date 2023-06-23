package userModel

import (
	"time"
)

// User model
type User struct {
	Id            string     `json:"id" bson:"_id,omitempty"`
	Email         *string    `json:"email" bson:"email,omitempty"`
	Locale        string     `json:"locale" bson:"locale,omitempty"`
	Fullname      string     `json:"fullname" bson:"fullname,omitempty"`
	Password      string     `json:"password" bson:"password,omitempty"`
	PasswordRetry string     `json:"retryPassword" bson:"-"`
	RecoveryToken string     `json:"recoveryToken" bson:"-"`
	Gender        string     `json:"gender" bson:"-"`
	Profiles      []Profile  `json:"profiles" bson:"profiles,omitempty"`
	CreatedAt     *time.Time `json:"created_at" bson:"created_at,omitempty"`
	UpdatedAt     *time.Time `json:"updated_at" bson:"updated_at,omitempty"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Email     *string   `json:"email"`
	Fullname  string    `json:"fullname"`
	Gender    string    `json:"gender"`
	Profiles  []Profile `json:"profiles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Profile struct {
	Type    string `json:"type"`
	Id      string `json:"id"`
	Link    string `json:"link"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}
