package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/irdaislakhuafa/go-sdk/codes"
	"github.com/irdaislakhuafa/go-sdk/errors"
)

func main() {
	uc := NewExampleUsecase()
	params := Entity{ID: uuid.NewString()}

	// Read Order: 1
	// We create data into db
	result, err := uc.Create(params)
	if err != nil {
		// Read Order: 4
		// more simple to identify, we can identify error by each context of error
		if code := errors.GetCode(err); code == codes.CodeInvalidValue {
			fmt.Printf("response 400 bad request, %v\n", err)
		} else if code >= codes.CodeSMTPStart && code <= codes.CodeSMTPEnd { // we can identify more than one context at same technology integration with simple code like this
			// additional case to handle SMTP error
			fmt.Println("respone bla bla bla")
		} else {
			fmt.Printf("response 500 internal server error: %v\n", err)
		}
	}
	fmt.Printf("result 200 success: %#v\n", result)

}

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

		// Imagine if email already exists
		if isEmailAlreadyExists := (time.Now().Unix()/2 == 0); isEmailAlreadyExists {
			return errors.NewWithCode(codes.CodeInvalidValue, "email %#v already exists", params.Email) // we can handle dynamic error messages
		}

		if params.Email == "" {
			return errors.NewWithCode(codes.CodeInvalidValue, "email is required and cannot be empty")
		}

		if params.Age <= 0 {
			return errors.NewWithCode(codes.CodeInvalidValue, "age is invalid")
		}

		if params.Name == "" {
			return errors.NewWithCode(codes.CodeInvalidValue, "name is required and cannot be empty")
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

	return result, nil
}

func (eu *exampleUsecaseImpl) GetList() ([]Entity, error) {
	results := []Entity{}
	return results, nil
}
