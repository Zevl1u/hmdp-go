package beans

type Result struct {
	Success bool        `json:"success,omitempty"`
	ErrMsg  string      `json:"err_msg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
