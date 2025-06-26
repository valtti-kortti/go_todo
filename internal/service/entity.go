package service

type TaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TaskUpdate struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
