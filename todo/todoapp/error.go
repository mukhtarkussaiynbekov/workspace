package todoapp

const (
	// ErrToDoItemNotFound when todo item is not found.
	ErrToDoItemNotFound = Error("todo item not found")
	// ErrToDoItemAlreadyExist when todo item already exist in the system.
	ErrToDoItemAlreadyExist = Error("todo item already exists")
)

// Error represents a Sampleapp error.
type Error string

// Error returns the error message.
func (e Error) Error() string {
	return string(e)
}
