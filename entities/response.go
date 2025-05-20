package entities

type Response struct {
	Status          string `json:"status"`
	StatusCode      int    `json:"statusCode,omitempty"`
	Message         string `json:"message,omitempty"`
	ErrorMessage    string `json:"errorMessage,omitempty"`
	ErrorCode       string `json:"errorCode,omitempty"`
	TransactionCode string `json:"transactionCode,omitempty"`
	Data            interface{}
}
