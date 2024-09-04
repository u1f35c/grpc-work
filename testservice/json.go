//go:generate protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative testservice.proto

package testservice

// JSONRequest represents the data transferred for a request in the JSON
// implementation of the Test Service
type JSONRequest struct {
	Action string
	Value  int
}

// JSONResponse represents the data transferred for a response in the JSON
// implementation of the Test Service
type JSONResponse struct {
	Value  int
}
