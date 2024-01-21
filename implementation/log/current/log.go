package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	uc := NewExampleUsecase()
	params := Entity{ID: uuid.NewString()}

	// We create data into db
	result, err := uc.Create(params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %#v\n", result)

	// We get list data from db
	results, err := uc.GetList()
	if err != nil {
		panic(err)
	}
	fmt.Printf("results: %#v\n", results)
}

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
	// Read Order: 1
	// Code to save the entity to database
	// ASK: but we want to know who the user created this data, how we can do it?
	return params, nil
}

func (eu *exampleUsecaseImpl) GetList() ([]Entity, error) {
	results := []Entity{}
	// Read Order: 2
	// Imagine when database performance is down
	// ASK: cliens will wait fora long time to complete their requests, how do we handle this?
	return results, nil
}
