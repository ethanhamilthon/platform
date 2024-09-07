package auth

import (
	"controller/internal/db"
	"encoding/json"
	"errors"
	"io"
)

type UserData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}


type AuthService struct {
	db *db.DB
}

func New(db *db.DB) *AuthService {
	return &AuthService{db: db}
}

type RegisterBody struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (a *AuthService) Register(r io.Reader) (string, error) {
	if a.isUserExists() {
		return "", errors.New("user already exists")
	}

	data := new(RegisterBody)
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		return "", err
	}

	if data.Password != data.ConfirmPassword {
		return "", errors.New("passwords do not match")
	}

	user := UserData{
		Username: data.Username,
		Password: data.Password,
	}

	body, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	err = a.db.Set("auth:user", body)
	if err != nil {
		return "", err
	}

	err = a.setUserExists()
	if err != nil {
		return "", err
	}
	token, err := a.GenerateJWT(1)
	if err != nil {
		return "", err
	}
	return token, nil
}

type LoginBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (a *AuthService) Login(r io.Reader) (string, error) {
	if !a.isUserExists() {
		return "", errors.New("user doesn't exist")
	}

	data := new(LoginBody)
	err := json.NewDecoder(r).Decode(data)
	if err != nil {
		return "", err
	}

	userData, err := a.db.Get("auth:user")
	if err != nil {
		return "", err
	}

	user := new(UserData)
	err = json.Unmarshal(userData, user)
	if err != nil {
		return "", err
	}

	if user.Username != data.Username || user.Password != data.Password {
		return "", errors.New("invalid credentials")
	}

	token, err := a.GenerateJWT(1)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *AuthService) isUserExists() bool {
	data, err := a.db.Get("is_user_exists")
	if err != nil {
		return false
	}
	return string(data) == "true"
}

func (a *AuthService) setUserExists() error {
	return a.db.Set("is_user_exists", []byte("true"))
}

type MeBody struct {
	Token string `json:"token"`
}

func (a *AuthService) Me(r io.Reader) (string, error) {
	body := new(MeBody)
	err := json.NewDecoder(r).Decode(body)
	if err != nil {
		return "", err
	}

	err = a.VerifyJWT(body.Token)
	if err != nil {
		return "", err
	}

	userData, err := a.db.Get("auth:user")
	if err != nil {
		return "", err
	}

	user := new(UserData)
	err = json.Unmarshal(userData, user)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

