package entity

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type Ltree string

func (Ltree) GormDataType() string {
	return "ltree"
}

func (l Ltree) Value() (driver.Value, error) {
	return string(l), nil
}

func (l *Ltree) Scan(value interface{}) error {
	v, ok := value.(string)
	if !ok {
		return fmt.Errorf("Ltree scan - cast %v to string err", value)
	}

	*l = Ltree(v)

	return nil
}

func (l Ltree) CastLquery(lquery string) string {
	return fmt.Sprintf("CAST ((%s) AS lquery)", lquery)
}

func AuthPathArrayToLtree(authPaths []string) string {
	var stmt = new(strings.Builder)
	stmt.WriteString("ARRAY[")
	for i := range authPaths {
		authPaths[i] = strings.ReplaceAll(authPaths[i][1:], "/", ".")
		authPaths[i] = strings.ReplaceAll(authPaths[i], "-", "_")
		if len(authPaths) == 0 {
			continue
		}
		stmt.WriteString("'")
		stmt.WriteString(authPaths[i])
		stmt.WriteString("'")
		if i < len(authPaths)-1 {
			stmt.WriteString(",")
		}
	}

	stmt.WriteString("]::lquery[]")

	return stmt.String()
}
