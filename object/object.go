package object

import "fmt"

const (
	INTEGER_OBJ = "integer"
	ERROR_OBJ = "error"
	FLOAT_OBJ = "float"
)

type ObjectType string

type Object interface {
	Type()		ObjectType
	Inspect()	string
}

type Integer struct{
	Value	int64
}


func (i *Integer) Type() ObjectType {return INTEGER_OBJ}
func (i *Integer) Inspect() string {return fmt.Sprintf("%d", i.Value)}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

type Float struct {
	Value	float64
}
func (f *Float) Type()ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string { return fmt.Sprintf("%v", f.Value) }
