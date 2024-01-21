package main

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uc := NewExampleUsecase()
	params := Entity{ID: uuid.NewString()}

	// Read Order: 1
	// We create data into db
	result, err := uc.Create(params)
	if err != nil {
		// Read Order: 4
		// if err is invalid field validation we want to return 400 (bad request to client) and we need to write this code
		switch err {
		case ErrEmailAlreadyExists, ErrEmailIsRequired, ErrAgeNotValid, ErrNameIsRequired: // can be long code
			fmt.Printf("response 400 bad request, %v\n", err)
		default:
			fmt.Printf("response 500 internal server error, %v", err)
		}
	}
	fmt.Printf("result 200 success: %#v\n", result)

}

// Imagine this error layer, developers need to write every annoying error statically (not dynamic)
var (
	ErrEmailIsRequired    = errors.New("email is required")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNameIsRequired     = errors.New("name is required")
	ErrAgeNotValid        = errors.New("age not valid")
)

// Imagine this is entity layer
type Entity struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	Age   int64  `json:"age,omitempty"`
}

// Imagine this usecase layer
type ExampleUsecase interface {
	Create(params Entity) (Entity, error)
	GetList() ([]Entity, error)
}

type exampleUsecaseImpl struct {
}

func NewExampleUsecase() ExampleUsecase {
	return &exampleUsecaseImpl{}
}

func (eu *exampleUsecaseImpl) Create(params Entity) (Entity, error) {
	// Read Order: 2
	// Imagine this is repository code to save data in db
	validate := func(params Entity) error {
		if params.Email == "" {
			return ErrEmailIsRequired
		}

		if params.Age <= 0 {
			return ErrAgeNotValid
		}

		if params.Name == "" {
			return ErrNameIsRequired
		}

		return nil
	}

	if err := validate(params); err != nil {
		return Entity{}, err
	}

	// Read Order: 3
	// Code to save the entity to database
	createInDB := func(params Entity) (Entity, error) {
		return params, nil
	}
	result, err := createInDB(params)
	if err != nil {
		return Entity{}, err
	}

	/*
		ASK:
		- how we can handle dynamic error message? some case we need to use dynamic error message, as example in validation, to make error message more simple like this "invalid validation field [email, age, name] cannot be empty or null" depends the invalid fields
	*/

	return result, nil
}

func (eu *exampleUsecaseImpl) GetList() ([]Entity, error) {
	results := []Entity{}
	return results, nil
}
