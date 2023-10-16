package models

type JSONResponse struct {
	Data      string `json:"data"`
	Sequence  uint64 `json:"sequence"`
	Timestamp int64  `json:"timestamp"`
}

type JSONRequest struct {
	Status int
	Data   string
}
