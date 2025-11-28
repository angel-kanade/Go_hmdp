package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type Response struct {
	BizCode int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"` // 返回数据，omitempty使其在为空时不输出
}

// WriteResponse used to write an error and JSON data into response.
func writeResponse(c *gin.Context, bizCode int, message string, data any) {
	coder, ok := codes[bizCode]
	if !ok {
		coder = codes[ErrUnknown]
	}

	if message != "" {
		coder.Message = message
	}

	if coder.HTTPStatus() != http.StatusOK {
		c.AbortWithStatusJSON(coder.HTTPStatus(), Response{
			BizCode: bizCode,
			Message: coder.Error(),
			Data:    data,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		BizCode: bizCode,
		Message: coder.Error(),
		Data:    data})
}

// response.go
func Success(c *gin.Context, data any) {
	writeResponse(c, ErrSuccess, "", data)
}

func SuccessWithMsg(c *gin.Context, message string, data any) {
	writeResponse(c, ErrSuccess, message, data)
}

func Error(c *gin.Context, bizCode int, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	writeResponse(c, bizCode, msg, nil)
}
func ErrorWithData(c *gin.Context, bizCode int, data any, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	writeResponse(c, bizCode, msg, data)
}

/*
*
在Web应用开发中，我们经常遇到这样的困扰：

	业务函数既要处理逻辑又要发送HTTP响应，职责不清
	错误日志记录分散在各处，难以统一管理
	错误信息对用户不友好，技术细节泄露给客户端
	原始错误信息丢失，排查问题困难
*/
type BusinessError struct {
	Code    int
	Message string
	Err     error `json:"-"` //原始错误
}

// 实现了Error() string，就是实现了error接口
func (e *BusinessError) Error() string {
	return e.Message
}

func (e *BusinessError) Unwrap() error {
	return e.Err
}

func NewBusinessError(code int, customMessage string) *BusinessError {
	if meta, exists := codes[code]; exists {
		message := customMessage
		if message == "" {
			message = meta.Message
		}
		return &BusinessError{
			Code:    code,
			Message: message,
		}
	}
	return &BusinessError{
		Code:    ErrUnknown,
		Message: "未知错误",
	}
}

// WrapBusinessError 包装原始错误
func WrapBusinessError(code int, originalErr error, customMessage string) *BusinessError {
	bizErr := NewBusinessError(code, customMessage)
	bizErr.Err = originalErr
	return bizErr
}

// response.go
func HandleBusinessError(c *gin.Context, err error) {
	var bizErr *BusinessError
	if errors.As(err, &bizErr) {
		//记录业务错误日志
		if bizErr.Err != nil {
			slog.Error("business error", "originalErr", bizErr.Err, "code", bizErr.Code, "message", bizErr.Message)
		}
		Error(c, bizErr.Code, bizErr.Message)
	} else {
		slog.Error("unknown error", "err", err)
		Error(c, ErrUnknown)
	}
}

// HandleBusinessErrorWithData 统一处理业务错误（可附带数据）
func HandleBusinessErrorWithData(c *gin.Context, err error, data any) {
	var bizErr *BusinessError
	if errors.As(err, &bizErr) {
		if bizErr.Err != nil {
			slog.Error("business error with data", "originalErr", bizErr.Err, "code", bizErr.Code, "message", bizErr.Message)
		}
		ErrorWithData(c, bizErr.Code, data, bizErr.Message)
	} else {
		slog.Error("unknown error with data", "err", err)
		ErrorWithData(c, ErrUnknown, data)
	}
}

// HandleBusinessResult 统一处理业务结果（错误或成功）
func HandleBusinessResult(c *gin.Context, err error, data any) {
	if err != nil {
		HandleBusinessError(c, err)
	} else {
		Success(c, data)
	}
}

func HandleBusinessResultWithErrorData(c *gin.Context, err error, data any) {
	if err != nil {
		HandleBusinessErrorWithData(c, err, data)
	} else {
		Success(c, data)
	}
}
