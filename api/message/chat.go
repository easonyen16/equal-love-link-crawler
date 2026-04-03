package message

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ChatMedia struct {
	ID                       int    `json:"id"`
	FileExtension            string `json:"fileExtension"`
	ContentType              string `json:"contentType"`
	URL                      string `json:"url"`
	CompressedURL            string `json:"compressedUrl"`
	ThumbnailURL             string `json:"thumbnailUrl"`
	CompressedThumbnailURL   string `json:"compressedThumbnailUrl"`
	ThumbnailContentType     string `json:"thumbnailContentType"`
	ThumbnailFileExtension   string `json:"thumbnailFileExtension"`
	DurationInSecs           int    `json:"durationInSecs"`
}

type ChatMessage struct {
	ID                   int         `json:"id"`
	TalkRoomID           int         `json:"talkRoomId"`
	PostedArtistID       int         `json:"postedArtistId"`
	PostedUserID         int         `json:"postedUserId"`
	TextContent          string      `json:"textContent"`
	Type                 string      `json:"type"`
	ReservationDate      int64       `json:"reservationDate"`
	PostedDate           int64       `json:"postedDate"`
	Status               string      `json:"status"`
	IsRead               bool        `json:"isRead"`
	FanLetterID          int         `json:"fanLetterId"`
	TargetUserID         int         `json:"targetUserId"`
	ChatMedia            []ChatMedia `json:"chatMedia"`
	IsMine               bool        `json:"isMine"`
	PostedUsername       string      `json:"postedUsername"`
	PostedUserProfileURL string      `json:"postedUserProfileUrl"`
	PostedUserIconURL    string      `json:"postedUserIconUrl"`
	IsFavorite           bool        `json:"isFavorite"`
}

type ChatPage struct {
	Messages       []ChatMessage
	NextPageID     int
	PreviousPageID int
}

type chatAPIResponse struct {
	Result         bool          `json:"result"`
	NextPageID     int           `json:"nextPageId"`
	PreviousPageID int           `json:"previousPageId"`
	Data           []ChatMessage `json:"data"`
}

func GetChat(accessToken string, talkRoomID, page, pageStartID int) (*ChatPage, error) {
	q := url.Values{}
	q.Set("page", fmt.Sprintf("%d", page))
	q.Set("pageSize", "100")
	q.Set("hasMedia", "false")
	q.Set("isFavorite", "false")
	q.Set("isSentFanLetter", "false")
	q.Set("dateSearchInSecs", "0")
	q.Set("orderBy", "1")
	if pageStartID != 0 {
		q.Set("pageStartId", fmt.Sprintf("%d", pageStartID))
	}

	u := url.URL{
		Scheme:   "https",
		Host:     AppDomain,
		Path:     fmt.Sprintf("/user/v2/chat/%d", talkRoomID),
		RawQuery: q.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", HeaderUserAgent)
	req.Header.Set("Accept-Language", HeaderAcceptLanguage)
	req.Header.Set("X-Request-Verification-Key", HeaderXRequestVerificationKey)
	req.Header.Set("X-Artist-Group-UUID", HeaderXArtistGroupUUID)
	req.Header.Set("X-Device-UUID", HeaderXDeviceUUID)
	req.Header.Set("Host", AppDomain)
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp chatAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Result {
		return nil, errors.New("get chat failed")
	}

	return &ChatPage{
		Messages:       apiResp.Data,
		NextPageID:     apiResp.NextPageID,
		PreviousPageID: apiResp.PreviousPageID,
	}, nil
}
