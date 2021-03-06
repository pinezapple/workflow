// Package controllers  use for decode -> function -> logs -> encode -> response
package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"workflow/heimdall/webserver/forms"

	"github.com/gin-gonic/gin"
)

const (
	defaultPage     string = "1"
	defaultPageSize string = "10"
)

var operators = [5]string{">=", "<=", ">", "<", "="}

// ResponseError return response error
func ResponseError(c *gin.Context, er error, statusCode int) {
	response := &forms.ResponseError{
		Message:    er.Error(),
		StatusCode: statusCode,
	}
	c.JSON(statusCode, response)
	c.Abort()
	return
}

func getFilterParam(c *gin.Context) (pageSize int, pageToken int, filter map[string][]string, err error) {
	query := c.Request.URL.Query()

	pageSizeStr := query.Get("page_size")
	if pageSizeStr == "" {
		pageSizeStr = defaultPageSize
	}
	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil {
		logger.Errorf("Error when parse page size: %s", err.Error())
		ResponseError(c, err, http.StatusBadGateway)
		return
	}

	pageTokenStr := query.Get("page_token")
	if pageTokenStr == "" {
		pageTokenStr = defaultPage
	}
	pageToken, err = strconv.Atoi(pageTokenStr)
	if err != nil {
		logger.Errorf("Error when parse page token: %s", err.Error())
		ResponseError(c, err, http.StatusBadRequest)
		return
	}

	filters := query.Get("filter")
	conditions := strings.Split(filters, ";")
	filter = make(map[string][]string, 0)

	for _, condition := range conditions {
		var filterPair []string

		for _, operator := range operators {
			if strings.Contains(condition, operator) {
				filterPair = strings.Split(condition, operator)
				filter[filterPair[0]] = []string{operator, filterPair[1]}
				break
			}
		}
	}
	return
}
