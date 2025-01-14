package response

import (
	"schedrestd/common"
	"schedrestd/common/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is the gin http response with JSON format
type Response struct {
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

// ResOK returns the response struct
func ResOK(data interface{}) Response {
	return Response{
		Message: common.SuccessMsg,
		Data:    data,
	}
}

// ResOKGinJson invoke gin JSON func directly
func ResOKGinJson(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Message: common.SuccessMsg,
		Data:    data,
	})
}

// ResErr handles the unexpected error
func ResErr(c *gin.Context, code int, err error) {
	logger.GetDefault().Errorf("%v", err.Error())
	c.JSON(code, Response{
		Message: err.Error(),
	})
}
