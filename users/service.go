package users

import (
	"errors"
	"alquran/helpers"

	"github.com/aws/aws-sdk-go/aws/session"
	"golang.org/x/crypto/bcrypt"
)

type Services interface {
	RegisterUser(s3 *session.Session, request RegisterUserRequest) (map[string]interface{}, error)
	Login(request LoginRequest) (User, error)
}

type services struct {
	repository Repository
}

func NewServices(repository Repository) *services {
	return &services{repository}
}

func (s *services) RegisterUser(s3 *session.Session, request RegisterUserRequest) (map[string]interface{}, error) {
	result := map[string]interface{}{}
	user := User{}
	user.Uuid = helpers.GenerateUUID()
	user.Username = request.Username
	user.Email = request.Email
	if request.Password != request.PasswordRetype {
		return result, errors.New("Password Not Match! Please Makesure Password Same.")
	}
	password, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return result, errors.New("Failed To Generate Password Hashing")
	}
	user.Password = string(password)

	newUser, err := s.repository.Save(user)
	if err != nil {
		err = helpers.HandleDuplicateKeyError(err)
		return result, err
	}

	// response result setelah di filter data apa aja yang akan di tampilkan
	result = map[string]interface{}{
		"uuid":     newUser.Uuid,
		"username": newUser.Username,
		"email":    newUser.Email,
	}
	return result, nil
}

func (s *services) Login(request LoginRequest) (User, error) {

	user, err := s.repository.FindByUsername(request.Username)
	if err != nil {
		return user, errors.New("User Not Found on Database, Please Check Your Credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return user, errors.New("Password Your Given Not Match! Password Wrong.")
	}

	return user, nil
}
