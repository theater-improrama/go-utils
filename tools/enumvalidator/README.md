
# enumvalidator

Generates validation methods for enum types automatically.

## Usage

In your module, include the following generation command (filename `mytype.go`):

```go
package mypackage
//go:generate go run github.com/theater-improrama/go-utils/tools/enumvalidator -type=MyType
```

This will generate a file called `mytype_enumvalidator.go` in the same directory as your source file.

