package main

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username    string `validate:"required,min=4,max=20"`
	Email       string `validate:"required,email"`
	Age         int    `validate:"gte=0,lte=130"`
	Gender      string `validate:"required,gender"`
	DateOfBirth string `validate:"required,dateformat"`
}

func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == "male" || gender == "female" || gender == "other"
}

func validateDateFormat(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	// Regex untuk format tanggal YYYY-MM-DD
	pattern := `^\d{4}-\d{2}-\d{2}$`
	matched, err := regexp.MatchString(pattern, date)
	if err != nil {
		return false
	}

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}

	return matched
}

func main() {
	user := User{
		Username:    "john_doe",
		Email:       "john@example.com",
		Age:         25,
		Gender:      "male",
		DateOfBirth: "1990-12-31",
	}

	validate := validator.New()
	validate.RegisterValidation("gender", validateGender)
	validate.RegisterValidation("dateformat", validateDateFormat)

	err := validate.Struct(user)
	if err != nil {
		// Validation failed
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(user).FieldByName(err.StructField())
			fieldName := field.Tag.Get("json")
			fmt.Printf("Validation failed on field %s with tag '%s'\n", fieldName, err.Tag())
		}
	} else {
		// Validation passed
		fmt.Println("Validation passed")
	}
}
