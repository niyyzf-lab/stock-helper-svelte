package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ErrorCode 错误代码类型
type ErrorCode int

const (
	// 执行引擎错误码
	ErrEngineAlreadyRunning ErrorCode = iota + 1000
	ErrEngineNotRunning
	ErrEnginePaused
	ErrEngineTimeout

	// 工作池错误码
	ErrWorkerPoolCreation
	ErrWorkerCreation
	ErrWorkerClosed

	// Lua相关错误码
	ErrLuaFunctionRegistration
	ErrLuaScriptExecution
	ErrLuaStateCreation

	// API相关错误码
	ErrAPIRequestFailed
	ErrAPIResponseInvalid

	// 数据相关错误码
	ErrInvalidStockData
	ErrInvalidKLineData

	// 配置相关错误码
	ErrInvalidConfig
	ErrInvalidStrategy
)

// ErrorLevel 错误级别
type ErrorLevel string

const (
	ErrorLevelInfo    ErrorLevel = "info"
	ErrorLevelWarning ErrorLevel = "warning"
	ErrorLevelError   ErrorLevel = "error"
	ErrorLevelFatal   ErrorLevel = "fatal"
)

// ErrorEvent 错误事件
type ErrorEvent struct {
	Code      ErrorCode  `json:"code"`
	Message   string     `json:"message"`
	Level     ErrorLevel `json:"level"`
	Time      time.Time  `json:"time"`
	Details   string     `json:"details,omitempty"`
	Component string     `json:"component"`
}

// EngineError 自定义错误类型
type EngineError struct {
	Code      ErrorCode
	Message   string
	Err       error
	Level     ErrorLevel
	Component string
}

// Error 实现error接口
func (e *EngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *EngineError) Unwrap() error {
	return e.Err
}

// ToEvent 转换为错误事件
func (e *EngineError) ToEvent() ErrorEvent {
	var details string
	if e.Err != nil {
		details = e.Err.Error()
	}

	return ErrorEvent{
		Code:      e.Code,
		Message:   e.Message,
		Level:     e.Level,
		Time:      time.Now(),
		Details:   details,
		Component: e.Component,
	}
}

// EmitError 发送错误事件到前端
func EmitError(ctx context.Context, err error) {
	if engineErr, ok := err.(*EngineError); ok {
		runtime.EventsEmit(ctx, "engine:error", engineErr.ToEvent())
	} else {
		// 包装普通错误
		engineErr := &EngineError{
			Code:      ErrEngineNotRunning,
			Message:   "Unknown error occurred",
			Err:       err,
			Level:     ErrorLevelError,
			Component: "Engine",
		}
		runtime.EventsEmit(ctx, "engine:error", engineErr.ToEvent())
	}
}

// NewEngineError 创建新的引擎错误
func NewEngineError(code ErrorCode, message string, err error) *EngineError {
	level := ErrorLevelError
	switch code {
	case ErrEngineAlreadyRunning, ErrEnginePaused:
		level = ErrorLevelInfo
	case ErrEngineTimeout, ErrWorkerClosed:
		level = ErrorLevelWarning
	case ErrLuaStateCreation, ErrWorkerPoolCreation:
		level = ErrorLevelFatal
	}

	return &EngineError{
		Code:      code,
		Message:   message,
		Err:       err,
		Level:     level,
		Component: "Engine",
	}
}

// IsEngineError 判断是否为引擎错误
func IsEngineError(err error) bool {
	_, ok := err.(*EngineError)
	return ok
}

// 错误创建辅助函数
func ErrAlreadyRunning() error {
	return NewEngineError(ErrEngineAlreadyRunning, "engine is already running", nil)
}

func ErrWorkerPoolFailed(err error) error {
	return NewEngineError(ErrWorkerPoolCreation, "failed to create worker pool", err)
}

func ErrLuaFuncRegFailed(err error) error {
	return NewEngineError(ErrLuaFunctionRegistration, "failed to register Lua functions", err)
}

func ErrLuaScriptFailed(err error) error {
	return NewEngineError(ErrLuaScriptExecution, "failed to execute Lua script", err)
}

func NewAPIRequestError(operation string, err error) error {
	return &EngineError{
		Code:      ErrAPIRequestFailed,
		Message:   fmt.Sprintf("API request failed: %s", operation),
		Err:       err,
		Level:     ErrorLevelWarning,
		Component: "API",
	}
}

func NewInvalidStockDataError(code string, err error) error {
	return &EngineError{
		Code:      ErrInvalidStockData,
		Message:   fmt.Sprintf("invalid stock data for %s", code),
		Err:       err,
		Level:     ErrorLevelWarning,
		Component: "Data",
	}
}

func NewInvalidConfigError(field string, err error) error {
	return &EngineError{
		Code:      ErrInvalidConfig,
		Message:   fmt.Sprintf("invalid configuration: %s", field),
		Err:       err,
		Level:     ErrorLevelError,
		Component: "Config",
	}
}
