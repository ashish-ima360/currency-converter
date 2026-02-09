package utils

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseIDParam(param string, c *gin.Context) (int, error) {
	idStr, ok := c.Params.Get(param)
	if !ok {
		return 0, errors.New("Invalid ID parameter")
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("ID must be an integer")
	}
	return id, nil
}
