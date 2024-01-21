# `Log`

Logging in Go serves the purpose or recording or printing messages to a log service or specific output. Logging has several important functions in software development, including.

- Troubleshooting

  Logging allows developer to records events and additional information while program is running. This help developer in troubleshooting and solving issues when application is in production environment.

- Performance analysis

  Logging allows developer to records performance or related information. Such as the time taken to execute a specific function or operation. This help developer to identifying areas that can can be optimized.

- Audit and Security

  Logging can be used to records security or related events or activities for audit purposes. This can including user access, configuration changes, or other security critical activities.

- Monitoring

  In production environment, logging is important part of monitoring system. The recorded information helps in understanding the overall health and performance of the system.

- Facilitating Maintenance

  Logs also provide insight during development process, helping developers to understand the program excution flow and capturing information along the way. This help developers for maintenance and ongoing development.

## Implementation

- Current

  Currently the Go application in SIESTA has implemented log, but in my personal opinion is not very efficient, because the developer needs to write which layer this code executed every time it logs and wherre this error occurs.

  Several weakness in current implementation of logs in SIESTA:

  1. The current log implementation cannot be used to monitor the execution time of each operation.
  2. Doesn't tell the developer the specific part of the code where the error occured. Developers needs to track the error manually by messages
  3. Cannot track errors on specific clients, whis will make difficult for developers to handle problems with specific operations if the application is being accessed by many clients.

  There following is an example of the log that siesta currently implements

  <details>
  <summary>Click here</summary>

  ```go
    package main

    import (
      "time"

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

      log.Infof("result: %+v", result)

      // We get list data from db
      results, err := uc.GetList()
      if err != nil {
        panic(err)
      }
      log.Infof("results: %+v\n", results)
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
      log.Info("STATE USECASE -> Create(), execute method to create data")

      // Imagine this is method from repository
      createToDB := func(params Entity) (Entity, error) {
        time.Sleep(time.Second * 5)
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

      log.Info("STATE USECASE -> Create(), returning response")
      return result, nil
    }

    func (eu *exampleUsecaseImpl) GetList() ([]Entity, error) {
      log.Info("STATE USECASE -> GetList(), get list data")

      // Imagine this is method to get list data from repository
      getListData := func() ([]Entity, error) {
        time.Sleep(time.Second * 5)
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
      log.Info("STATE USECASE -> GetList(), returning response")
      return results, nil
    }


  ```

  </details>

  The result of this implementation looks like this. Run with command `go run implementation/log/current/log.go`
  <details>
  <summary>Click here</summary>

  ```bash
    2024/01/22 03:58:53.023223 log.go:52: [Info] STATE USECASE -> Create(), execute method to create data
    2024/01/22 03:58:58.027485 log.go:76: [Info] STATE USECASE -> Create(), returning response
    2024/01/22 03:58:58.027641 log.go:21: [Info] result: {ID:922149bd-d080-4f20-a837-04a26d214533 Email: Name: Age:0}
    2024/01/22 03:58:58.027676 log.go:81: [Info] STATE USECASE -> GetList(), get list data
    2024/01/22 03:59:03.031628 log.go:101: [Info] STATE USECASE -> GetList(), returning response
    2024/01/22 03:59:03.031703 log.go:28: [Info] results: []
  ```

  </details>

- Recommendation

  The following is the logs implementation that i recommend to answer the problem above. This implementation will be easier if we implement `context.Context`.

  Here i use my [`personal library`](https://github.com/irdaislakhuafa/go-sdk.git) as an example, i hope the SIESTA team can make the implementation better.
  <details>
  <summary>Click here</summary>

  ```go
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


  ```

  </details>

  My recommendation implementation will looks like this. Run with `go run implementation/log/recommendation/log.go`

  <details>
  <summary>Click here</summary>

  ```json
    {
    "level": "info",
    "request_id": "fe6363a7-527c-4e36-88ed-6292e84ed5fe",
    "service_version": "",
    "time_elapsed": "0ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:31+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:99",
    "message": "create data"
  }
  {
    "level": "info",
    "request_id": "fe6363a7-527c-4e36-88ed-6292e84ed5fe",
    "service_version": "",
    "time_elapsed": "5003ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:36+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:123",
    "message": "returning response"
  }
  {
    "level": "info",
    "request_id": "fe6363a7-527c-4e36-88ed-6292e84ed5fe",
    "service_version": "",
    "time_elapsed": "5003ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:36+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:31",
    "message": "result: {ID:b6f15f97-9537-47d3-a801-fd1962a443ae Email: Name: Age:0 CreatedBy:}\n"
  }
  {
    "level": "info",
    "request_id": "c3061909-ffb1-4912-aa60-147ed435039b",
    "service_version": "",
    "time_elapsed": "0ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:36+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:129",
    "message": "get list data"
  }
  {
    "level": "info",
    "request_id": "c3061909-ffb1-4912-aa60-147ed435039b",
    "service_version": "",
    "time_elapsed": "5003ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:41+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:150",
    "message": "returning response"
  }
  {
    "level": "info",
    "request_id": "c3061909-ffb1-4912-aa60-147ed435039b",
    "service_version": "",
    "time_elapsed": "5003ms",
    "user_agent": "",
    "time": "2024-01-22T03:49:41+07:00",
    "caller": "/media/Projects/Companies/SIESTA/Repositories/golang-boilerplate/implementation/log/recommendation/log.go:41",
    "message": "results: []\n"
  }
  ```

  </details>

  - `request_id`: We can filter each log by specific user request by `request_id`, this can be different value depending on the case
  - `time_elapsed`: We can monitor the time required for each operation
  - `caller`: We can monitor via this field to find out which line of code that need optimization
  - `level`: We can filter specific level of logs for troubleshooting (ex, `error`/`debug`/`info`/etc)

  I recommend the JSON format to make it easier to integrate with monitoring applications such as Grafana or others.
