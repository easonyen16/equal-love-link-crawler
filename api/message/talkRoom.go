package message

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Artist struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	RomajiName       string `json:"romajiName"`
	ImageURL         string `json:"imageUrl"`
	IconURL          string `json:"iconUrl"`
	OriginalImageURL string `json:"originalImageUrl"`
	OriginalIconURL  string `json:"originalIconUrl"`
}

type TalkRoom struct {
	ID                            int      `json:"id"`
	OwnedArtistGroupID            int      `json:"ownedArtistGroupId"`
	Name                          string   `json:"name"`
	RomajiName                    string   `json:"romajiName"`
	NeedPurchaseFlag              bool     `json:"needPurchaseFlag"`
	UUID                          string   `json:"uuid"`
	LastMessageDate               int64    `json:"lastMessageDate"`
	ImageURL                      string   `json:"imageUrl"`
	Unread                        int      `json:"unread"`
	LastSeen                      int64    `json:"lastSeen"`
	Artists                       []Artist `json:"artists"`
	IsAccessible                  bool     `json:"isAccessible"`
	IOSSubscriptionProductSKU     string   `json:"iosSubscriptionProductSku"`
	AndroidSubscriptionProductSKU string   `json:"androidSubscriptionProductSku"`
}

type TalkRoomsResponse struct {
	TalkRooms                    []TalkRoom `json:"talkRooms"`
	TotalUnreadCount             int        `json:"totalUnreadCount"`
	TotalUnreadNotificationCount int        `json:"totalUnreadNotificationCount"`
	NotArrivedNotifications      []any      `json:"notArrivedNotifications"`
}

type talkRoomsAPIResponse struct {
	Result bool              `json:"result"`
	Data   TalkRoomsResponse `json:"data"`
}

func (c *Client) GetTalkRooms(accessToken string) ([]TalkRoom, error) {
	q := url.Values{}
	q.Set("page", "1")

	u := url.URL{
		Scheme:   "https",
		Host:     c.cfg.appDomain,
		Path:     "/user/v2/talk-room",
		RawQuery: q.Encode(),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}
	c.setHeaders(req, c.cfg.appDomain, "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var apiResp talkRoomsAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	if !apiResp.Result {
		return nil, errors.New("get talk rooms failed")
	}

	return apiResp.Data.TalkRooms, nil
}
