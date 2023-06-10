package user

import (
	"errors"
	"net/http"

	"github.com/handymesh/handy_authService/db/mongodb"
	"github.com/handymesh/handy_authService/models/user"
	"github.com/handymesh/handy_authService/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckUniqueUser(w http.ResponseWriter, user userModel.User) bool {
	count, err := mongodb.Session.Database("auth").Collection(userModel.CollectionUser).CountDocuments(nil, bson.M{"mail": user.Email})
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`))
		return true
	}

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New(`{"mail": "need unique mail"}`))
		return true
	}

	return false
}
