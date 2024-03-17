package task

import (
	"encoding/json"

	"github.com/OliviaDilan/async_arc/pkg/contract"
)

type TaskCreatedV1 struct {
	contract.Metadata `json:"metadata"`

	TaskID int `json:"task_id"`
}

func NewTaskCreatedV1(taskID int) *TaskCreatedV1 {
	return &TaskCreatedV1{
		TaskID: taskID,
		Metadata: contract.NewMetadataV1(),
	}
}

func (m TaskCreatedV1) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *TaskCreatedV1) Unmarshal(body []byte) error {
	return json.Unmarshal(body, m)
}

//task completed contract

type TaskCompletedV1 struct {
	contract.Metadata `json:"metadata"`

	TaskID int `json:"task_id"`
	UserID string `json:"user_id"`
}

func NewTaskCompletedV1(taskID int, userID string) *TaskCompletedV1 {
	return &TaskCompletedV1{
		TaskID: taskID,
		UserID: userID,
		Metadata: contract.NewMetadataV1(),
	}
}


func (m TaskCompletedV1) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *TaskCompletedV1) Unmarshal(body []byte) error {
	return json.Unmarshal(body, m)
}