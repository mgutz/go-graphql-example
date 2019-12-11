package scalars

import (
	"encoding/json"
)

// Any is a  GraphQL scalar to represent any valid JSON value  which can be
// a string, int or map
type Any struct {
	json.RawMessage
}

// NewAny creates an instance of Any.
func NewAny(v interface{}) *Any {
	a := &Any{}
	a.UnmarshalGraphQL(v)
	return a
}

// ImplementsGraphQLType maps this custom Go type
// to the graphql scalar type in the schema.
func (Any) ImplementsGraphQLType(name string) bool {
	return name == "Any"
}

// UnmarshalGraphQL is a custom unmarshaler for Any
//
// This function will be called whenever you use the
// Any scalar as an input
func (a *Any) UnmarshalGraphQL(input interface{}) error {
	switch input := input.(type) {
	case []byte:
		a.RawMessage = input
	case json.RawMessage:
		a.RawMessage = input
	default:
		var err error
		a.RawMessage, err = json.Marshal(input)
		return err
	}

	return nil
}

// MarshalJSON is a custom marshaler for Any
//
// This function will be called whenever you
// query for fields that use the Time type
func (a Any) MarshalJSON() ([]byte, error) {
	return a.RawMessage, nil
}
