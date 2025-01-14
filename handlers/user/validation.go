package user

import (
	"errors"
	"fmt"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/handymesh/hyshAuthService/db/mongodb"
	userModel "github.com/handymesh/hyshAuthService/models/user"
	"github.com/handymesh/hyshAuthService/utils"
)

func CheckUniqueUser(w http.ResponseWriter, user userModel.User) bool {
	count, err := mongodb.Session.Database("auth").Collection(userModel.CollectionUser).CountDocuments(nil, bson.M{"mail": user.Email})
	if err != nil {
		utils.Error(w, errors.New(`"`+err.Error()+`"`), http.StatusBadRequest)
		return true
	}

	fmt.Printf("User %v count %v", bson.M{"mail": user.Email}, count)

	if count > 0 {
		w.WriteHeader(http.StatusBadRequest)
		utils.Error(w, errors.New(`{"mail": "need unique mail"}`), http.StatusBadRequest)
		return true
	}

	return false
}
