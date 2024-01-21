# `Error`

Error handling in Golang is important to ensuring that applications run and can handle unexpected situations.
Here is some key functions of error handling in Golang:

- Separation of Normal Flow and Error Flow: error handling allows a clear separation between the normal execution path of the program and the handling of error conditions. This enables easy identification of sections of the program dealing with normal situations and those handling errors.
- Indicate failure conditions: error handling allows a function to indicate that something has gone wrong or a failure has occured. For example, function in Golang returned `error` value and make decision based on the result, this helps prevent oversight in handling error conditions.
- Error handling at appropriate level: Golang supports multiple return values, allowing functions to return values along with an error. This allows the caller to easily check and handle any potential errors.

## Implementation

- Current

  Currently, SIESTA has implemented quite good error handling, but it cannot be used in some conditions and requires code that long enough to identify errors.

  The following are some of the weaknesses of error handling in SIESTA:

  - Cannot be used to handle errors with dynamic error messages
  - Requires a long error checking code for the same error context (example error at validation)

  There following is an example of the error handling that siesta currently implements.

  <details>
  <summary>Click here</summary>

  ```go
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


  ```

  </details>

- Recommendation

  The following is the error handling implementation that i recommend to solve the problem above.

  Here i use my [`personal library`](https://github.com/irdaislakhuafa/go-sdk.git) as an example, i hope the SIESTA team can make the implementation better.

  <details>
  <summary>Click here</summary>

  ```go
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

  ```

  </details>

  OK, maybe that's all my recommendation for now, I hope the Siesta team can consider it :I
