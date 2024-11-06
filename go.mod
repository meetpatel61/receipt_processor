// go.mod - Defines the module's dependencies and configuration

// Specifies the name of the module. The module name is typically the project's
// path or repository location. Here, it's named `receipt_processor`.
module receipt_processor

// Specifies the version of the Go language used for this module.
// In this case, `go 1.23.2` specifies that the project is using Go version 1.23.2.
go 1.23.2

// Lists the dependencies required by this module. Each dependency specifies a package
// and the version required. In this case, the project depends on `github.com/google/uuid`,
// a package used for generating unique IDs. Version `v1.6.0` is specified as the required version.
require github.com/google/uuid v1.6.0