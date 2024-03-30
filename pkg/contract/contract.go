package contract

import (
	"errors"
	"time"
)

type Contract interface {
	contract()
	Marshal()([]byte, error)
	Unmarshal(body []byte)error
}

type Metadata struct {
	Timestamp time.Time   `json:"timestamp"`
	Version   Version `json:"version"`
}

type Version string

const (
	VersionV1 Version = "v1"
)

func NewMetadataV1() Metadata {
	return Metadata{
		Timestamp: time.Now(),
		Version:   VersionV1,
	}
}

func (m Metadata) contract() {
}

var ErrContractUnimplemented = errors.New("unimplemented contract")

func (m Metadata) Marshal() ([]byte, error) {
	return nil, ErrContractUnimplemented
}

func (m *Metadata) Unmarshal(body []byte) error {
	return ErrContractUnimplemented
}


