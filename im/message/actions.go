package message

type Action string

const (
	ActionMessage            Action = "message"
	ActionGroupMessage              = "message.group"
	ActionChatMessage               = "message.chat"
	ActionChatMessageRecall         = "message.chat.recall"
	ActionGroupMessageRecall        = "message.group.recall"
	// ActionChatMessageRetry 消息重发, 服务器未ack
	ActionChatMessageRetry = "message.chat.retry"
	// ActionChatMessageResend 消息重发, 服务器已ack, 接收方未ack
	ActionChatMessageResend = "message.chat.resend"
	ActionCSMessage         = "message.cs"
	ActionMessageFailed     = "message.failed.send"

	ActionNotifyNeedAuth      = "notify.auth"
	ActionNotifyKickOut       = "notify.kickout"
	ActionNotifyNewContact    = "notify.contact"
	ActionNotifyGroup         = "notify.group"
	ActionNotifyAccountLogin  = "notify.login"
	ActionNotifyAccountLogout = "notify.logout"
	ActionNotifyError         = "notify.error"

	ActionAckRequest  = "ack.request"
	ActionAckGroupMsg = "ack.group.msg"
	ActionAckMessage  = "ack.message"
	ActionAckNotify   = "ack.notify"

	ActionApi       = "api"
	ActionHeartbeat = "heartbeat"
	ActionApiFailed = "api.failed"
)
