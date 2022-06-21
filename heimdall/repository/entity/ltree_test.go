package entity_test

import (
	"fmt"
	"testing"

	"github.com/vfluxus/heimdall/repository/entity"
)

func TestAuthPathArrayToLtree(t *testing.T) {
	var (
		authPaths = []string{"/a/b/c-1-2-3", "/c/d/e"}
	)

	stmt := entity.AuthPathArrayToLtree(authPaths)
	fmt.Println(stmt)
}
