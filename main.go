package main

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Username    *string `validate:"required,base64,min=4,max=9"`
	Email       string  `validate:"required,email"`
	Age         int     `validate:"gte=0,lte=130"`
	Gender      string  `validate:"required,gender"`
	DateOfBirth string  `validate:"required,dateformat"`
	ScoringType string  `validate:"-"`
	Score       string  `validate:"x=aiforesee"`
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

func ValidateDecodeBase64(fl validator.FieldLevel) bool {
	date := fl.Field().String()

	decodeString, err := base64.StdEncoding.DecodeString(date)
	if err != nil {
		return false
	}

	decoded := string(decodeString)
	if !fl.Field().CanSet() { //pengecekan nilai pointer
		return false
	}
	fl.Field().SetString(string(decoded)) //nilai struct yg akan ditimpa harus berupa pointer atau akan error
	return true
}

func validateRequiredByScoringType(fl validator.FieldLevel) bool {
	// Mendapatkan nilai pada properti yang di validasi
	value := fl.Field().String()

	fmt.Println("value :", value)
	// Mendapatkan nilai properti ScoringType dari struct parent
	scoringType := fl.Parent().FieldByName("ScoringType").String()

	fmt.Println("scoringType :", scoringType)

	// Mengambil nilai parameter tag
	param := fl.Param()
	fmt.Println("param :", param)

	switch scoringType {
	case param:
		if value == "" {
			return false
		}
	}

	return true
}

func main() {

	//username := "john_doe"
	username := "am9obl9kb2U="
	email := "john@example.com"
	age := 25
	gender := "male"
	dateOfBirth := "1990-12-31"
	ScoringType := "aiforesee"
	scoring := "a"

	user := User{
		Username:    &username,
		Email:       email,
		Age:         age,
		Gender:      gender,
		DateOfBirth: dateOfBirth,
		ScoringType: ScoringType,
		Score:       scoring,
	}

	validate := validator.New()
	validate.RegisterValidation("gender", validateGender)
	validate.RegisterValidation("dateformat", validateDateFormat)
	validate.RegisterValidation("base64", ValidateDecodeBase64)
	validate.RegisterValidation("x", validateRequiredByScoringType)

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

	fmt.Println(user)
	// fmt.Println(*user.Username)
}
