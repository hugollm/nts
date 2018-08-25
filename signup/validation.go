package signup

import (
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/mail"
	"strings"
)

type SignupInput struct {
	Email    string
	Password string
}

func validateInput(body io.Reader) (SignupInput, map[string]error) {
	errors := make(map[string]error)
	input, jsonErr := validateJson(body)
	if jsonErr != nil {
		errors["json"] = jsonErr
		return SignupInput{}, errors
	}
	vEmail, emailErr := validateEmail(input.Email)
	if emailErr != nil {
		errors["email"] = emailErr
	}
	vPassword, passwordErr := validatePassword(input.Password)
	if passwordErr != nil {
		errors["password"] = passwordErr
	}
	return SignupInput{vEmail, vPassword}, errors
}

func validateJson(body io.Reader) (SignupInput, error) {
	var input SignupInput
	jsonErr := json.NewDecoder(body).Decode(&input)
	if jsonErr != nil {
		return input, errors.New("Invalid JSON.")
	}
	return input, nil
}

func validateEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return email, errors.New("Email is required.")
	}
	parsed, parseErr := mail.ParseAddress(email)
	if parseErr != nil {
		return email, errors.New("Invalid email.")
	}
	email = parsed.Address
	if userWithEmailExists(email) {
		return email, errors.New("Email is already taken.")
	}
	return email, nil
}

func validatePassword(password string) (string, error) {
	if password == "" {
		return password, errors.New("Password is required.")
	}
	if len(password) < 8 {
		return password, errors.New("Password must be at least 8 characters long.")
	}
	return hashPassword(password), nil
}

func hashPassword(password string) string {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)
}
