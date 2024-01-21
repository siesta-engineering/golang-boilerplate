package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func main() {
	// Imagine this is middleware layer
	ctx := ExampleMiddleware()

	uc := NewExampleUsecase()

	params := Entity{ID: uuid.NewString()}
	// We create data into db with context.Context
	result, err := uc.Create(ctx, params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("result: %#v\n", result)

	// We get list data from db
	results, err := uc.GetList(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Printf("results: %#v\n", results)
}

// Imagine this is middleware layer
const CtxKeyUserID = "user_id"

func ExampleMiddleware(args ...any) context.Context {
	ctx := context.Background()

	// Read Order: 2
	// For usecase like "Create()" we can handle auth here and passed user_id into context value
	// Do some logic to authenticate user
	ctx = context.WithValue(ctx, CtxKeyUserID, uuid.NewString())

	// Read Order: 4
	// For usecase like "GetList" we can setting timeout ot operation in context, as example we set maximum timeout is 5 seconds
	// We can save timeout configuration in ".env" or other resource to make it more flexible to customize
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel() // this method will be executed if operation take more than 5 seconds to complete and this method will stopped all operation that uses this context

	// Imagine this context is used in http.Request
	return ctx
}

// Imagine this is entity layer
type Entity struct {
	ID        string `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Name      string `json:"name,omitempty"`
	Age       int64  `json:"age,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}

// Imagine this usecase layer
type ExampleUsecase interface {
	Create(ctx context.Context, params Entity) (Entity, error)
	GetList(ctx context.Context) ([]Entity, error)
}

type exampleUsecaseImpl struct {
}

func NewExampleUsecase() ExampleUsecase {
	return &exampleUsecaseImpl{}
}

func (eu *exampleUsecaseImpl) Create(ctx context.Context, params Entity) (Entity, error) {
	// Read Order: 1
	// Code to save the entity to database
	// ASK: but we want to know who the user created this data, how we can do it?

	// With context we can get addition information without modifying the function or interfaces like below
	params.CreatedBy = ctx.Value(CtxKeyUserID).(string)

	// Save to DB

	return params, nil
}

func (eu *exampleUsecaseImpl) GetList(ctx context.Context) ([]Entity, error) {
	results := []Entity{}

	// Read Order: 3
	// Imagine when database performance is down
	// ASK: cliens will wait fora long time to complete their requests, how do we handle this?

	/*
		With context, we don't need to do anything here if we have used context in all methods, because in "context.Context" which is used in "GetList()" methods the operation will stop automatically if the previously determined timout has passed (at Read Order: 4).
		"eu.repository.GetList(ctx)"
	*/

	return results, nil
}
