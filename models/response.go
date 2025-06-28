package models

type ResponseMessage struct {
	Status  int         `json:"status"`
	Message interface{} `json:"message"`
}

type SMSResponse struct {
	Errors       []string `json:"errors"`
	ResponseCode int      `json:"response_code"`
	Status       string   `json:"status"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
