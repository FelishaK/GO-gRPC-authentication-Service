package authdto

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuthedUserResponse struct {
	FirstName        string
	LastName         string
	Email            string
	RoleName         string
	RegistrationDate primitive.DateTime
}
