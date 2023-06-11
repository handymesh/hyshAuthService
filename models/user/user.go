package userModel

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/handymesh/handy_authService/db/mongodb"
	"github.com/handymesh/handy_authService/utils/crypto"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

func List() (error, []User) {
	var users []User
	//err := mongodb.Session.Database("auth").Collection(CollectionUser).Find(nil).Sort("-updated_on").All(&users)
	//cursor, err := mongodb.Session.Database("auth").Collection(CollectionUser).Find(nil, users)
	//if err != nil {
	//	return err, users
	//}

	return nil, users
}

func Find(user User) (*[]User, error) {
	var users *[]User
	res := mongodb.Session.Database("auth").Collection(CollectionUser).FindOne(nil, user)
	res.Decode(users)
	if res.Err() != nil {
		return nil, res.Err()
	}

	return users, nil
}

func FindOne(user User) (*User, error) {
	var result *User
	err := mongodb.Session.Database("auth").Collection(CollectionUser).FindOne(nil, user).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func FindCount(user User) (int64, error) {
	count, err := mongodb.Session.Database("auth").Collection(CollectionUser).CountDocuments(nil, user)

	if err != nil {
		return count, err
	}

	return count, nil
}

func Add(user User) (error, User) {
	var checkUser User
	checkUser.Email = user.Email

	result, err := FindCount(checkUser)
	if err != nil {
		return err, user
	}

	if result > 0 {
		err := errors.New(`{"mail":"need unique mail`)
		fmt.Printf("our err is %v /n", user)
		return err, user
	}

	user.Password, _ = crypto.HashPassword(user.Password)
	fmt.Printf("User password %v", user.Password)
	time := time.Now()
	user.CreatedAt = &time
	user.UpdatedAt = &time

	res, err := mongodb.Session.Database("auth").Collection(CollectionUser).InsertOne(nil, user)
	if err != nil {
		return errors.New(`{"mail":"need unique mail"}`), user
	}

	user.Id = res.InsertedID.(primitive.ObjectID).Hex()

	return nil, user
}

func Update(user *User) (*User, error) {
	UpdatedAt := time.Now()
	user.UpdatedAt = &UpdatedAt
	user.Email = nil // prohibit changing address

	filter := bson.D{{"_id", user.Id}}
	_, err := mongodb.Session.Database("auth").Collection(CollectionUser).UpdateOne(nil, filter, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Delete(userId string) (int64, error) {
	filter := bson.D{{"_id", userId}}
	res, err := mongodb.Session.Database("auth").Collection(CollectionUser).DeleteOne(nil, filter)

	return res.DeletedCount, err
}
