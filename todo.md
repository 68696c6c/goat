
# v0.2.0
- add linting
- replace pagination with hypertext linking?
  x R&D how to implement HAL
  x need to R&D how to manage pagination params with the new hypertext pagination links
- ID to Id, DB to Db, etc
x R&D generic repos & controllers
- validation error responder (RespondValidationError takes a map instead of an error so it doesn't match the ErrorResponder interface)
- more utils
x upgrade to gorm v2
x upgrade validator to v10 (github.com/go-playground/validator)
  - no longer using validator in favor of repo methods like Create and Update
x context
x ENV var: currently only used to assume HTTP_DEBUG based on env, remove in favor of just setting HTTP_DEBUG
x export goatDB.ConnectionConfig



# current
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
- filter sorting works, but should work with json names instead of db column names
  - or should it? camel case makes sense fo js clients, snake => camel is much simpler than camel => snake in code gen, so maybe snake makes more sense system-wide...
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
    x with child
    x with parent
    - with both relations
  - hard delete model:
    - with no relations
    - with child
    - with parent
    - with both relations

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
