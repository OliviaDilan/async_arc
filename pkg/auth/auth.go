package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	DecodeToken(token string) (*User, error)
	GetUsers() ([]*User, error)
}

type authClient struct {
	client *http.Client
	host string
}

func NewClient(hostName, port string) Client {
	return &authClient{
		host: hostName+":"+port,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   5 * time.Second,
		},
	}
}

func (c *authClient) DecodeToken(token string) (*User, error) {

	req := decodeTokenRequest{
		Token: token,
	}

	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(req); err != nil {
		return nil, err
	}
	
	resp, err := c.client.Post(c.host+"/decode_token", "application/json", buf)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var res decodeTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &User{
		Username: res.Username,
		Role:     res.Role,
	}, nil
}

func (c *authClient) GetUsers() ([]*User, error) {

	resp, err := c.client.Get(c.host + "/users")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var res getUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Users, nil
}