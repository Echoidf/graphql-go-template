package response

type Response struct {
	Data    any    `json:"data,omitempty"`
	Msg     string `json:"msg,omitempty"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

func Success(data any, msg string) Response {
	return Response{
		Data:    data,
		Success: true,
		Msg:     msg,
	}
}
