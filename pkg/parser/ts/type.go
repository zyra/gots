package ts

type BasicType string

const (
	BasicTypeString  BasicType = "string"
	BasicTypeNumber  BasicType = "number"
	BasicTypeBoolean BasicType = "boolean"
	BasicTypeAny     BasicType = "any"
)

type Type struct {
	Name  string
	Array bool
}
