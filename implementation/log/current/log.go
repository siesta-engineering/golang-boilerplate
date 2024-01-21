package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2/log"

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
	log.Infof("STATE USECASE -> Create(), execute method to create data")

	// Imagine this is method from repository
	createToDB := func(params Entity) (Entity, error) {
		return params, nil
	}

	// Read Order: 1
	// Code to save the entity to database
	result, err := createToDB(params)
	if err != nil {
		// Developer need to write where this error occurred manually for each logs
		// If the code in the "Create()" method layer of this usecase has reached many lines, it will take time to find out where the error line occurs.
		log.Errorf("STATE USECASE -> Create(), err: %v", err)
		return Entity{}, err
	}

	/*
		ASK:
		- why don't we create a log interface that tells us specifically on which line of this error occurred?
		- how do we track the logs for specific user_id = XXX that we use to track errors if many clients access them simultaneously?
	*/
	return result, nil
}

func (eu *exampleUsecaseImpl) GetList() ([]Entity, error) {
	log.Errorf("STATE USECASE -> GetList(), get list data")

	// Imagine this is method to get list data from repository
	getListData := func() ([]Entity, error) {
		return []Entity{}, nil
	}

	// Read Order: 2
	// Imagine when database performance is down and the server takes several minutes to retrieve data because it needs optimization when querying the database afte a lot data has been stored.
	results, err := getListData()
	if err != nil {
		log.Errorf("STATE USECASE -> GetList(), %v", err)
		return nil, err
	}

	/*
		ASK:
		- how do we know which line of code takes the slowest time to execute if we intergate this method with other technology? how do we know how long this operation takes to execute?
	*/

	return results, nil
}
