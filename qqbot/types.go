package qqbot

type SendMsgReply struct {
	ID              string `json:"id"`
	ChannelID       string `json:"channel_id"`
	GuildID         string `json:"guild_id"`
	Content         string `json:"content"`
	Timestamp       string `json:"timestamp"`
	TTS             bool   `json:"tts"`
	MentionEveryone bool   `json:"mention_everyone"`
	Author          Author `json:"author"`
	Pinned          bool   `json:"pinned"`
	Type            int64  `json:"type"`
	Flags           int64  `json:"flags"`
	SeqInChannel    string `json:"seq_in_channel"`
	ErrorMessage    string `json:"message"`
	ErrorCode       int    `json:"code"`
}

type Author struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
}
type Message struct {
	Content        string `json:"content"`
	ReplyMessageID string `json:"msg_id"`
}

func NewMessage(content string, replyto string) *Message {
	return &Message{Content: content, ReplyMessageID: replyto}
}
