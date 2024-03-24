package domain_model

import (
	"fmt"
	"strconv"
)

type ParseError struct {
    Type string
    Value string
}

func (e ParseError) Error() string {
    return fmt.Sprintf("Failed to parse '%v' as %v", e.Value, e.Type)
}

type Id int

func ParseId(idstr string) (Id, error) {
    id, err := strconv.ParseUint(idstr, 10, 32)
    if err != nil {
        return 0, ParseError{ Type: "id", Value: idstr }
    }
    return Id(id), nil
}
