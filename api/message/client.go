package message

type Platform int

const (
	PlatformEqualLove Platform = iota
	PlatformNotEqualMe
)

type platformConfig struct {
	authDomain string
	appDomain  string
	headers    map[string]string
}

var platformConfigs = map[Platform]platformConfig{
	PlatformEqualLove: {
		authDomain: "api.entertainment-platform-auth.cosm.jp",
		appDomain:  "v3.api.equal-love.link.cosm.jp",
		headers: map[string]string{
			"User-Agent":                 "io.cosm.fc.user.equal.love/1.2.0/Android/16/SM-A217F",
			"Accept-Language":            "ja",
			"X-Request-Verification-Key": "rvk_78BBFC93-045B-4B88-A671-D2BB5B533164",
			"Content-Type":               "application/json",
			"X-Artist-Group-UUID":        "1d85e1ca-d845-4594-a8d8-ad36db293f01",
			"X-Device-UUID":              "android_BP4A.251205.006",
		},
	},
	PlatformNotEqualMe: {
		authDomain: "api.entertainment-platform-auth.cosm.jp",
		appDomain:  "v3.api.not-equal-me.link.cosm.jp",
		headers: map[string]string{
			"User-Agent":                 "io.cosm.fc.user.not.equal.me/1.3.0/Android/16/SM-A217F",
			"Accept-Language":            "ja",
			"X-Request-Verification-Key": "rvk_78BBFC93-045B-4B88-A671-D2BB5B533164",
			"Content-Type":               "application/json",
			"X-Artist-Group-UUID":        "aed48b34-b17b-4082-a42e-2bdc0d3022ff",
			"X-Device-UUID":              "android_BP4A.251205.006",
		},
	},
}

type Client struct {
	cfg platformConfig
}

func NewClient(platform Platform) *Client {
	return &Client{cfg: platformConfigs[platform]}
}
