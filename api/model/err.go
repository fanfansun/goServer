package model

type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

type ErrResponse struct {
	HttpSc int // 错误码
	Error  Err
}

// 定义错误

var (
	// 请求错误
	ErrorRequestBodyParseFailed = ErrResponse{HttpSc: 400, Error: Err{Error: "request - fail", ErrorCode: "001"}}
	// 未验证用户错误
	ErrorNotAuthUser = ErrResponse{HttpSc: 401, Error: Err{Error: "AuthUser - fail", ErrorCode: "002"}}
	// DB数据库错误
	ErrorDBError = ErrResponse{HttpSc: 500, Error: Err{Error: "DB ops failed", ErrorCode: "003"}}
	// 网络服务错误
	ErrorInternalFaults = ErrResponse{HttpSc: 500, Error: Err{Error: "Internal service error", ErrorCode: "004"}}
)
