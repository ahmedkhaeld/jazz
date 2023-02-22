package jazz

import "fmt"

//////////ErrPath//////////

type ErrPath struct {
	Path    string
	Message string
	Cause   error
}

var (
	ErrPathFail     = &ErrPath{Path: "", Message: "fail creating path", Cause: nil}
	ErrPathNotFound = &ErrPath{Path: "", Message: "path not found", Cause: nil}
	ErrFileFail     = &ErrPath{Path: "", Message: "fail create the required file", Cause: nil}
	ErrEnvFileFail  = &ErrPath{Path: "", Message: "fail create .env", Cause: nil}
)

func (e *ErrPath) Error() string {
	return fmt.Sprintf("path: %q: Message: %s: Cause: %v", e.Path, e.Message, e.Cause)
}

func (e *ErrPath) Is(target error) bool {
	t, ok := target.(*ErrPath)
	if !ok {
		return false
	}
	return t.Path == e.Path
}

func (e *ErrPath) Unwrap() error {
	return e.Cause
}

//////////ErrDB//////////

type ErrDB struct {
	Database string
	Message  string
	Cause    error
}

var (
	ErrDBNotOpen      = &ErrDB{Database: "", Message: " fail creating database ", Cause: nil}
	ErrDBNotConnected = &ErrDB{Database: "", Message: "fail connecting database", Cause: nil}
	ErrDBDial         = &ErrDB{Database: "", Message: "fail dialing redis", Cause: nil}
	ErrDBPing         = &ErrDB{Database: "", Message: "fail pinging redis", Cause: nil}
)

func (e *ErrDB) Error() string {
	return fmt.Sprintf("Database: %q: Message: %s: Cause: %v", e.Database, e.Message, e.Cause)
}

func (e *ErrDB) Is(target error) bool {
	t, ok := target.(*ErrDB)
	if !ok {
		return false
	}
	return t.Database == e.Database
}

func (e *ErrDB) Unwrap() error {
	return e.Cause
}
