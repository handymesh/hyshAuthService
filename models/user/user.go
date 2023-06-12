package userModel

import (
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/handymesh/hyshAuthService/db/mongodb"
	"github.com/handymesh/hyshAuthService/utils/crypto"
)

const (
	// CollectionUser holds the name of the articles collection
	CollectionUser = "users"
)

func List() (error, []User) {
	var users []User
	opts := options.Find().SetSort(bson.D{{"created_at", 1}})
	cursor, err := mongodb.Session.Database("auth").Collection(CollectionUser).Find(nil, bson.D{}, opts)
	if err != nil {
		return err, nil
	}

	defer cursor.Close(nil)

	// Process the users retrieved from the database as needed
	if err = cursor.All(nil, &users); err != nil {
		log.Fatal(err)
	}

	for i := range users {
		users[i].Password = ""
		users[i].PasswordRetry = ""
		users[i].RecoveryToken = ""
	}

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
	user.Id = "" // prohibit changing address

	opts := options.Update().SetUpsert(false)
	filter := User{Email: user.Email}
	fmt.Println("User is", user)
	update := bson.D{
		{"$set", user},
	}
	result, err := mongodb.Session.Database("auth").Collection(CollectionUser).UpdateOne(nil, filter, update, opts)
	if err != nil {
		return nil, err
	}
	fmt.Println("result", result)

	if result.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if result.UpsertedCount != 0 {
		fmt.Printf("inserted a new document with ID %v\n", result.UpsertedID)
	}

	return user, nil
}

func Delete(userId string) (int64, error) {
	filter := bson.D{{"_id", userId}}
	res, err := mongodb.Session.Database("auth").Collection(CollectionUser).DeleteOne(nil, filter)

	return res.DeletedCount, err
}
