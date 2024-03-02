package handler

type createTaskRequest struct {
	Title string `json:"title"`
}

type taskResponse struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Assignee    string  `json:"assignee"`
	Status      string  `json:"status"`
}

type createTaskResponse struct {
	Task taskResponse `json:"task"`
}

type getTasksResponse struct {
	Tasks []taskResponse `json:"tasks"`
}

type getTasksByAssigneeResponse struct {
	Tasks []taskResponse `json:"tasks"`
}

type closeTaskRequest struct {
	TaskID int `json:"id"`
}