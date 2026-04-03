package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	AuthDomain = "api.entertainment-platform-auth.cosm.jp"
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

func Login(username, password string) (*LoginResponse, error) {
	u := url.URL{
		Scheme: "https",
		Host:   AuthDomain,
		Path:   "/login",
	}

	body, err := json.Marshal(loginRequest{
		Username:   username,
		Password:   password,
		DeviceUUID: HeaderXDeviceUUID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", HeaderUserAgent)
	req.Header.Set("Accept-Language", HeaderAcceptLanguage)
	req.Header.Set("Content-Type", HeaderContentType)
	req.Header.Set("X-Request-Verification-Key", HeaderXRequestVerificationKey)
	req.Header.Set("X-Artist-Group-UUID", HeaderXArtistGroupUUID)
	req.Header.Set("X-Device-UUID", HeaderXDeviceUUID)
	req.Header.Set("Host", AuthDomain)
	req.Header.Set("Authorization", "Bearer")

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

func RefreshToken(accessToken, refreshToken string) (*LoginResponse, error) {
	u := url.URL{
		Scheme: "https",
		Host:   AuthDomain,
		Path:   "/refresh",
	}

	body, err := json.Marshal(refreshRequest{
		RefreshToken: refreshToken,
		DeviceUUID:   HeaderXDeviceUUID,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", HeaderUserAgent)
	req.Header.Set("Accept-Language", HeaderAcceptLanguage)
	req.Header.Set("Content-Type", HeaderContentType)
	req.Header.Set("X-Request-Verification-Key", HeaderXRequestVerificationKey)
	req.Header.Set("X-Artist-Group-UUID", HeaderXArtistGroupUUID)
	req.Header.Set("X-Device-UUID", HeaderXDeviceUUID)
	req.Header.Set("Host", AuthDomain)
	req.Header.Set("Authorization", "Bearer "+accessToken)

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
