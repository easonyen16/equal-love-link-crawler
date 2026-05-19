package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	UUID         string `json:"uuid"`
}

type loginRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	DeviceUUID string `json:"deviceUuid"`
}

type loginAPIResponse struct {
	Result bool          `json:"result"`
	Data   LoginResponse `json:"data"`
}

func (c *Client) Login(username, password string) (*LoginResponse, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.cfg.authDomain,
		Path:   "/login",
	}

	body, err := json.Marshal(loginRequest{
		Username:   username,
		Password:   password,
		DeviceUUID: c.cfg.headers["X-Device-UUID"],
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req, c.cfg.authDomain, "Bearer")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp loginAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Result {
		return nil, fmt.Errorf("login failed")
	}

	return &apiResp.Data, nil
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
	DeviceUUID   string `json:"deviceUuid"`
}

func (c *Client) RefreshToken(accessToken, refreshToken string) (*LoginResponse, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.cfg.authDomain,
		Path:   "/refresh",
	}

	body, err := json.Marshal(refreshRequest{
		RefreshToken: refreshToken,
		DeviceUUID:   c.cfg.headers["X-Device-UUID"],
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	c.setHeaders(req, c.cfg.authDomain, "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp loginAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Result {
		return nil, fmt.Errorf("refresh token failed")
	}

	return &apiResp.Data, nil
}

func (c *Client) setHeaders(req *http.Request, host, authorization string) {
	for k, v := range c.cfg.headers {
		req.Header.Set(k, v)
	}
	req.Header.Set("Host", host)
	req.Header.Set("Authorization", authorization)
}
