package main

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username string `tson:"username" validate:"required,min=4,max=20"`
	Email    string `validate:"required,email"`
	Age      int    `validate:"gte=0,lte=130"`
	Gender   string `validate:"required,gender"`
}

func validateGender(fl validator.FieldLevel) bool {
	gender := fl.Field().String()
	return gender == "male" || gender == "female" || gender == "other"
}

func main() {
	user := User{
		Username: "joh",
		Email:    "john@example.com",
		Age:      25,
		Gender:   "male",
	}

	validate := validator.New()
	validate.RegisterValidation("gender", validateGender)

	err := validate.Struct(user)
	if err != nil {
		// Validation failed
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := reflect.TypeOf(user).FieldByName(err.StructField())
			fieldName := field.Tag.Get("tson")
			fmt.Printf("Validation failed on field %s with tag '%s'\n", fieldName, err.Tag())
		}
	} else {
		// Validation passed
		fmt.Println("Validation passed")
	}
}
