package mongoDB

import (
	"context"
	"errors"
	"fmt"
	"searchandfind/pkg/models"

	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserModel struct {
	DB *mongo.Client
}

func (m *UserModel) Insert(u *models.User, password string) error {
	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	u.HashedPassword = passwordHashed

	collection := m.DB.Database("searchandfind").Collection("users")
	_, err = collection.InsertOne(context.TODO(), u)
	if err != nil {
		return err
	}
	return nil

}

func (m *UserModel) Authenticate(email, password string) (interface{}, error) {
	// filter := bson.M{"email": email}
	var user models.User
	fmt.Println(email, password)

	collection := m.DB.Database("searchandfind").Collection("users")
	fmt.Println(collection)
	err := collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)
	fmt.Println(err)
	if err != nil {
		fmt.Println("Coś jest nie tak!")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("Password is not ok!")
		return nil, models.ErrInvalidCredentials
	} else if err != nil {
		return nil, err
	}
	fmt.Println(user.ID.Hex())
	return user.ID.Hex(), nil
}

func (m *UserModel) Get(u string) (*models.User, error) {
	s := &models.User{}

	objID, _ := primitive.ObjectIDFromHex(u)
	collection := m.DB.Database("searchandfind").Collection("users")

	err := collection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&s)

	if err == models.ErrNoRecord {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *UserModel) Exists(email string) error {
	filter := bson.D{{"email", email}}
	var user models.User
	// fmt.Println("OK")
	// fmt.Println(filter)

	collection := m.DB.Database("searchandfind").Collection("users")
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	// fmt.Println(err)
	// fmt.Println(user.Name)
	if err != nil {
		return nil
	}
	return errors.New("Is not empty")

}
