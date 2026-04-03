package message

const (
	AppDomain = "v3.api.equal-love.link.cosm.jp"
)

type MessageApp interface {
	GetChat(accessToken string, talkRoomID, page, pageStartID int) (*ChatPage, error)
	GetTalkRooms(accessToken string) ([]TalkRoom, error)
}
