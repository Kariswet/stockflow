package model

type MetadataResponse struct {
	Status        bool   `json:"status"`
	Message       string `json:"message"`
	TimeExecution string `json:"time_execution"`
}

type Response struct {
	Metadata MetadataResponse `json:"metadata"`
	Data     interface{}      `json:"data"`
}

type PaginatedResponse struct {
	Data  interface{} `json:"data"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
}
