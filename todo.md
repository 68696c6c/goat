
# v0.2.0
- Id to ID, Db to DB, Url to URL, Http to HTTP, etc
x use snake case in query filtering, json names, etc!!!
  x even better, just use single word params.
x replace pagination with hypertext linking?
  x R&D how to implement HAL
  x need to R&D how to manage pagination params with the new hypertext pagination links
x R&D generic repos & controllers
x validation error responder (RespondValidationError takes a map instead of an error so it doesn't match the ErrorResponder interface)
x more utils
x upgrade to gorm v2
x upgrade validator to v10 (github.com/go-playground/validator)
  x no longer using validator in favor of repo methods like Create and Update
x context
x ENV var: currently only used to assume HTTP_DEBUG based on env, remove in favor of just setting HTTP_DEBUG
x export goatDB.ConnectionConfig

# v0.2.1
- add linting
- add http_request and http_testing tests
- add more examples tests
- R&D snake case vs camel case in json names and url params
  - url params can be case sensitive so snake is probably better (for example, query params in nextjs on vercel are case sensitive)

# old
x remove the http service struct meta service in favor of binding + repo Create/Update
x pure function controllers
  - take gin.Context as an arg, but return values and let goat handle responses
  - this should make it simpler to test controllers
    - not really... still need to pass a gin.Context to the controllers and since gin.CreateTestContext requires an http.ResponseWriter, it's about the same
x FILTER ENDPOINTS AND PAGINATION
  - filter users by name
  - filter users by current user type: supers -> all, admins, users -> own org
  - search users by name, email, org
x context
x gorm preload all: https://gorm.io/docs/preload.html#Preload-All
x router with built in link handling
x FILTERS AND PAGINATION 2.0
  - ~~there might be a simpler way to do this using gorm scopes: https://gorm.io/docs/scopes.html~~
    - not really... can't figure out how to dynamically group conditions :(
  x instead, use an updated, simplified query.Builder from goat-rnd
    x more tests!!!
    x examples of general purpose querying
    x examples of filtering w/ pagination
    x query builder + pagination from gin
- CAN WE BUILD USER PERMISSIONS INTO THE REPOS INSTEAD OF CONTROLLERS???
- apparently, the go convention for getters/setters is Thing() for get, SetThing() for set; update query package etc to work that way: https://go.dev/doc/effective_go#Getters
- filter sorting works, but should work with json names instead of db column names
  - or should it? camel case makes sense for js clients but, snake => camel is much simpler than camel => snake in code gen, so maybe snake makes more sense system-wide...
x maybe don't error out if base url isn't set during init?  would like to use goat for cli tools where that isn't necessary...
  - Init() only returns an error if base url is set and invalid
  - base url is not necessary unless InitRouter() is called, which will return an error if it isn't set
  - InitRouter() also now accepts an optional base url parameter to make things more flexible
- finish cleaning up outstanding todos!
- more example repos w/ 100% test coverage
  - crud
    x examples
    - tests
  - create-only
  - read-only
  - update-only
- delete-only
- more example models
  - soft delete model:
    - with no relations
    x with child (organizations)
    x with parent (users)
    - with both relations
  - hard delete model:
    - with no relations
    - with child
    - with parent
    - with both relations
x routes dilemma
  - SOLUTION:
    - do integration testing at the http package level with routes, middleware etc
    - controller tests are optional but should focus solely on unit testing the controller actions in isolation without routes or middleware, possibly with mocked repos?
      - TODO: add examples of these tests
        - set up Oauth authorization, add real user credentials and real auth middlewares
        - use the real auth stuff in the http test, use mocks in controller tests
          - authentication:
            - use client credentials to get a jwt for a user
  - current state:
    - controller tests need to use the same routes as the cmd/server in order to be useful
    - the routes config needs to be able reference app, middlewares, and controllers
    - middlewares should never need to reference controllers, but controllers may need to reference middleware-related things
      - for example, if a middleware sets the currently authenticated user, a controller may need to reference that user
  - thoughts:
    - controllers shouldn't really have logic that _needs_ tests; they should only be calling functions from other packages that have coverage
    - the thing to test is that things are combined correctly, i.e., a request + authentication = expected response
    - therefore, the tests could be done at the http level (this is the reasoning behind the current project structure where controllers are part of an app resource and routes are defined outside the controllers package)
    - however, controllers may still need to reference middleware-related things and middlewares aren't resource-related...
      - BUT, anything that a middleware sets should be set _in the context_ which controller handlers receive as an argument, so this is probably fine?
      - the only problem is the controller needs to know what the context key is in order access it, so that key needs to be defined in a place that is accessible to both the router/middlewares AND controllers... the most logical place might be in the controllers package
    - basically, a request is the entire handler chain (middlewares and controller action), so any request-based testing should be testing the entire chain
      - an actual controller unit test would use mocked context (hard) and repos (easy)!





eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ6WW5vQWItZEZfRl9lVGdHaWJfMjlaYkZEUEYyNU96Snl6SzNVQ3ZvT0ZZPSIsImV4cCI6MTY4MTE1NzQ4NX0.zquAH67Nh9jvOkNC7jlfOr2Ac-O_qWLcsVajBnQrsyLqQUzlhnGydL7bpea0ixyTwNTiL5eX5RvRhqWGU7yFvQ
# notes
- NewX vs MakeX conventions: 
  - use NewX for "constructors", i.e. functions that return interfaces, primitives (like errors), and empty structs or for functions that set struct values to their zeroed values or using arguments
  - use MakeX for "factories", i.e. functions that return structs with values set by the function itself instead of passed args, like functions for creating test mocks
  - in other languages, like C#, a constructor is called when the `new` keyword is used, while a factory is a method on a class that returns an instance
  - Notes on new() vs make() built-ins
    - the Go new() built-in allocates memory, but creates zeroed values (i.e. a nil array)
    - the Go make() built-in _initializes_ slices, maps, and channels, creating ready-to-use values (i.e. an array with a length and capacity)


# bugs
x pagination query params are being url formatted
x updatedAt appears to be set, but since it is gorm:"-", it isn't read

github.com/gin-contrib/zap.CustomRecoveryWithZap.func1.1
/app/src/vendor/github.com/gin-contrib/zap/zap.go:142
runtime.gopanic
/usr/local/go/src/runtime/panic.go:884
encoding/json.(*encodeState).marshal.func1
/usr/local/go/src/encoding/json/encode.go:327
runtime.gopanic
/usr/local/go/src/runtime/panic.go:884
runtime.panicmem
/usr/local/go/src/runtime/panic.go:260
runtime.sigpanic
/usr/local/go/src/runtime/signal_unix.go:835
github.com/68696c6c/example/app/models.(*Organization).getEmbedded
/app/src/app/models/organization.go:46
github.com/68696c6c/example/app/models.(*Organization).MarshalJSON
/app/src/app/models/organization.go:60
encoding/json.marshalerEncoder
/usr/local/go/src/encoding/json/encode.go:478
encoding/json.arrayEncoder.encode
/usr/local/go/src/encoding/json/encode.go:915
encoding/json.sliceEncoder.encode
/usr/local/go/src/encoding/json/encode.go:888
encoding/json.structEncoder.encode
/usr/local/go/src/encoding/json/encode.go:760
encoding/json.(*encodeState).reflectValue
/usr/local/go/src/encoding/json/encode.go:359
encoding/json.(*encodeState).marshal
/usr/local/go/src/encoding/json/encode.go:331
encoding/json.Marshal
/usr/local/go/src/encoding/json/encode.go:160
github.com/gin-gonic/gin/render.WriteJSON
/app/src/vendor/github.com/gin-gonic/gin/render/json.go:71
github.com/gin-gonic/gin/render.JSON.Render
/app/src/vendor/github.com/gin-gonic/gin/render/json.go:57
github.com/gin-gonic/gin.(*Context).Render
/app/src/vendor/github.com/gin-gonic/gin/context.go:910
github.com/gin-gonic/gin.(*Context).JSON
/app/src/vendor/github.com/gin-gonic/gin/context.go:953
github.com/gin-gonic/gin.(*Context).AbortWithStatusJSON
/app/src/vendor/github.com/gin-gonic/gin/context.go:204
github.com/68696c6c/goat.RespondOk
/app/src/vendor/github.com/68696c6c/goat/response.go:30
github.com/68696c6c/goat/controller.HandleList[...]
/app/src/vendor/github.com/68696c6c/goat/controller/controller.go:51
github.com/68696c6c/example/app/controllers.organizations.List
/app/src/app/controllers/organizations.go:29
github.com/gin-gonic/gin.(*Context).Next
/app/src/vendor/github.com/gin-gonic/gin/context.go:173
github.com/gin-contrib/zap.CustomRecoveryWithZap.func1
/app/src/vendor/github.com/gin-contrib/zap/zap.go:158
github.com/gin-gonic/gin.(*Context).Next
/app/src/vendor/github.com/gin-gonic/gin/context.go:173
github.com/gin-contrib/zap.GinzapWithConfig.func1
/app/src/vendor/github.com/gin-contrib/zap/zap.go:55
github.com/gin-gonic/gin.(*Context).Next
/app/src/vendor/github.com/gin-gonic/gin/context.go:173
github.com/gin-gonic/gin.(*Engine).handleHTTPRequest
/app/src/vendor/github.com/gin-gonic/gin/gin.go:616
github.com/gin-gonic/gin.(*Engine).ServeHTTP
/app/src/vendor/github.com/gin-gonic/gin/gin.go:572
github.com/68696c6c/goat.(*HandlerTest).Send
/app/src/vendor/github.com/68696c6c/goat/http_testing.go:40
github.com/68696c6c/example/app/controllers.Test_OrganizationsController_List
/app/src/app/controllers/organizations_test.go:14
testing.tRunner
/usr/local/go/src/testing/testing.go:1446
