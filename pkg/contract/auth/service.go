package auth

import (
	"encoding/json"

	"github.com/OliviaDilan/async_arc/pkg/contract"
)

type UserCreatedV1 struct {
	contract.Metadata `json:"metadata"`

	Username string `json:"username"`
}

func NewUserCreatedV1(username string) *UserCreatedV1 {
	return &UserCreatedV1{
		Username: username,
		Metadata: contract.NewMetadataV1(),
	}
}


func (m UserCreatedV1) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *UserCreatedV1) Unmarshal(body []byte) error {
	return json.Unmarshal(body, m)
}