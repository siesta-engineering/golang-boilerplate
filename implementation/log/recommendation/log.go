package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/irdaislakhuafa/go-sdk/appcontext"
	"github.com/irdaislakhuafa/go-sdk/log"
	"github.com/rs/zerolog"
)

func main() {
	// Read Order: 1
	// We can initialize log interface at first program running
	// We can save level configuration at ".env" or other location that make we (developer) easier to customize
	log := log.Init(log.Config{Level: zerolog.LevelDebugValue})

	// Read Order: 2
	// Imagine this is middleware layer
	ctx := ExampleMiddleware(log)

	uc := NewExampleUsecase(log)

	// Read Order: 6
	params := Entity{ID: uuid.NewString()}
	// We create data into db with context.Context
	result, err := uc.Create(ctx, params)
	if err != nil {
		panic(err)
	}

	log.Info(ctx, fmt.Sprintf("result: %+v\n", result))

	// Read Order: 9
	// In real implementation, different method is in different request, so they will have different context automatically
	ctx = ExampleMiddleware(log)
	// We get list data from db
	results, err := uc.GetList(ctx)
	if err != nil {
		panic(err)
	}
	log.Info(ctx, fmt.Sprintf("results: %+v\n", results))

}

// Read Order: 3
// Imagine this is middleware layer
func ExampleMiddleware(log log.Interface, args ...any) context.Context {
	// Read Order: 4
	// Imagine this is context from r.Context()
	ctx := context.Background()

	authenticateUser := func(ctx context.Context) (Entity, error) {
		return Entity{ID: uuid.NewString(), Name: "Irda Islakhu Afa"}, nil
	}

	// Do some logic to authenticate user
	user, err := authenticateUser(ctx)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed authenticate user, %v", err))
		return ctx
	}

	// Read Order: 5
	// we can set user_id in context as request_id to track specific user request if many clients access same resource simultaneously
	ctx = appcontext.SetRequestID(ctx, user.ID)
	// We can set start time execution of each operation in middleware to know which operation takes slowest time to be executed, time will added automatically to logs based on time value in context
	ctx = appcontext.SetRequestStartTime(ctx, time.Now())

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
	log log.Interface
}

func NewExampleUsecase(log log.Interface) ExampleUsecase {
	return &exampleUsecaseImpl{
		log: log,
	}
}

func (eu *exampleUsecaseImpl) Create(ctx context.Context, params Entity) (Entity, error) {
	// Read Order: 7
	// Developers do not need to write code location or other information every time they use logs
	eu.log.Info(ctx, "create data")

	// Imagine this is method from repository
	createToDB := func(params Entity) (Entity, error) {
		time.Sleep(time.Second * 5)
		return params, nil
	}

	// Code to save the entity to database
	result, err := createToDB(params)
	if err != nil {
		// Read Order: 8
		eu.log.Error(ctx, err)
		return Entity{}, err
	}

	eu.log.Info(ctx, "returning response")

	return result, nil
}

func (eu *exampleUsecaseImpl) GetList(ctx context.Context) ([]Entity, error) {
	// Read Order: 10
	// Developers do not need to write code location or other information every time they use logs
	eu.log.Info(ctx, "get list data")

	// Imagine this is method to get list data from repository
	getListData := func() ([]Entity, error) {
		time.Sleep(time.Second * 5)
		return []Entity{}, nil
	}

	// Read Order: 2
	// Imagine when database performance is down and the server takes several minutes to retrieve data because it needs optimization when querying the database afte a lot data has been stored.
	results, err := getListData()
	if err != nil {
		// Read Order: 11
		eu.log.Error(ctx, err)
		return nil, err
	}

	eu.log.Info(ctx, "returning response")
	return results, nil
}
