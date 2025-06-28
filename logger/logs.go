package logger

import (
	"caregiver-shift-tracker/models"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// response makes the response with payload as json format
func response(c *gin.Context, status int, payload []byte) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json; charset=UTF-8")
	c.Data(status, "application/json; charset=UTF-8", payload)
}

// RespondError makes the error response with payload as json format
func RespondError(c *gin.Context, code int, message string) {
	m := models.ErrorMessage{Error: message}
	res := models.ResponseMessage{Status: code, Message: m}
	msg, _ := json.Marshal(res)
	response(c, code, msg)
}

// Response makes the response with payload as json format
func Response(c *gin.Context, code int, message models.ResponseMessage) {
	msg, _ := json.Marshal(message)
	response(c, code, msg)
}

// RespondJSON makes the response with payload as json format
func RespondJSON(c *gin.Context, code int, message interface{}) {
	res := models.ResponseMessage{Status: code, Message: message}
	msg, _ := json.Marshal(res)
	response(c, code, msg)
}

// RespondRaw makes the response with payload as json format
func RespondRaw(c *gin.Context, code int, message interface{}) {
	msg, _ := json.Marshal(message)
	response(c, code, msg)
}

// RespondString makes the response with payload as string format
func RespondString(c *gin.Context, code int, message string) {
	c.String(code, message)
}
