package dbnames

import (
	"fmt"
)

const (
	allFields = "*"
)

type (
	SchemaName = string
	TableName  = string
	ViewName   = string
	FieldName  = string
)

type tableT struct {
	TName TableName
}

func (m *tableT) AllFields() FieldName {
	return fmt.Sprintf("%s.%s", m.TName, allFields)
}
