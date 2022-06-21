// libs provide shared function that shared by multiple packages
package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

// PrintJSONFormat print data with indent json format
func PrintJSONFormat(msg string, data interface{}) {
	if printThis, err := json.MarshalIndent(data, "", "    "); err == nil {
		fmt.Println(msg)
		fmt.Println(string(printThis))
		return
	} else {
		fmt.Println("---- Can not unmarshal to print, error: ", err)
		return
	}
}

func GetJwtToken(ctx *gin.Context) string {
	authzHeader := ctx.GetHeader("Authorization")
	return authzHeader[len("Bearer"):]
}
