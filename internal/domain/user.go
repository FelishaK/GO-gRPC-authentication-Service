package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id           primitive.ObjectID `bson:"_id,omitempty"`
	FirstName    string	`bson:"first_name"`
	LastName     string	`bson:"last_name"`
	PasswordHash []byte	`bson:"password_hash"`
	Email        string	`bson:"email"`
	RegistrationDate primitive.DateTime `bson:"registration_data,omitempty"`
	DateOfBirth primitive.DateTime `bson:"date_of_birth,omitempty"`
	Gender       string `bson:"gender,omitempty"`
	Role 				 *string `bson:"role,omitempty"`
}