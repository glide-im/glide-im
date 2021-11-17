package msgdao

var instance MsgDao

func init() {
	instance = impl{
		groupMsgDao: groupMsgDao{},
		chatMsgDao:  chatMsgDao{},
		cacheDao:    cacheDao{},
		commonDao:   commonDao{},
	}
}

type GroupMsgDao interface {
	GetGroupMsgSeq(gid int64) (int64, error)
	UpdateGroupMsgSeq(gid int64, seq int64) error

	GetGroupMessage(mid int64) (*GroupMessage, error)
	GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error)

	AddGroupMessage(from, to int64, cliMsgID string, type_ int, content string) (*GroupMessage, error)
	UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAg int64, lastMsgSeq int64) error
	GetGroupMessageState(gid int64) (*GroupMessageState, error)

	UpdateGroupMemberMsgState(gid int64, uid int64, lastAck int64, lastAckSeq int64) error
	GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error)
}

type ChatMsgDao interface {
	GetChatMessage(mid int64) (*ChatMessage, error)
	AddOrUpdateChatMessage(message *ChatMessage) (bool, error)

	GetChatMessageSeqAfter(uid int64, seqAfter int64) ([]*ChatMessage, error)
	GetChatMessageSeqSpan(uid int64, seq int64) (int, error)

	AddOfflineMessage(uid int64, mid int64) error
	GetOfflineMessage(uid int64) ([]*OfflineMessage, error)
	DelOfflineMessage(uid int64, mid []int64) error
}

type CacheDao interface {
	// GetUserMsgSeq 获取用户全当前局消息 Seq
	GetUserMsgSeq(uid int64) (int64, error)
	// GetIncrUserMsgSeq 返回用户全局消息递增seq, 保证递增, 尽量保持连续, 不保证一定连续
	GetIncrUserMsgSeq(uid int64) (int64, error)
}

type CommonDao interface {
	GetMessageID() (int64, error)
}

type MsgDao interface {
	ChatMsgDao
	GroupMsgDao
	CacheDao
	CommonDao
}

type impl struct {
	groupMsgDao
	chatMsgDao
	cacheDao
	commonDao
}

/////////////////

func GetUserMsgSeq(uid int64) (int64, error) {
	return instance.GetUserMsgSeq(uid)
}

func GetIncrUserMsgSeq(uid int64) (int64, error) {
	return instance.GetIncrUserMsgSeq(uid)
}

/////////////////

func GetMessageID() (int64, error) {
	return instance.GetMessageID()
}

/////////////////

func GetGroupMsgSeq(gid int64) (int64, error) {
	return instance.GetGroupMsgSeq(gid)
}
func UpdateGroupMsgSeq(gid int64, seq int64) error {
	return instance.UpdateGroupMsgSeq(gid, seq)
}
func GetGroupMessage(mid int64) (*GroupMessage, error) {
	return instance.GetGroupMessage(mid)
}
func GetGroupMessageSeqAfter(gid int64, seqAfter int64) ([]*GroupMessage, error) {
	return instance.GetGroupMessageSeqAfter(gid, seqAfter)
}
func AddGroupMessage(from, to int64, cliMsgID string, type_ int, content string) (*GroupMessage, error) {
	return instance.AddGroupMessage(from, to, cliMsgID, type_, content)
}
func UpdateGroupMessageState(gid int64, lastMID int64, lastMsgAg int64, lastMsgSeq int64) error {
	return instance.UpdateGroupMessageState(gid, lastMID, lastMsgAg, lastMsgSeq)
}
func GetGroupMessageState(gid int64) (*GroupMessageState, error) {
	return instance.GetGroupMessageState(gid)
}
func UpdateGroupMemberMsgState(gid int64, uid int64, lastAck int64, lastAckSeq int64) error {
	return instance.UpdateGroupMemberMsgState(gid, uid, lastAck, lastAckSeq)
}
func GetGroupMemberMsgState(gid int64, uid int64) (*GroupMemberMsgState, error) {
	return instance.GetGroupMemberMsgState(gid, uid)
}

///////////////////////////////////////

func GetChatMessage(mid int64) (*ChatMessage, error) {
	return instance.GetChatMessage(mid)
}
func AddChatMessage(message *ChatMessage) (bool, error) {
	return instance.AddOrUpdateChatMessage(message)
}
func GetChatMessageSeqAfter(uid int64, seqAfter int64) ([]*ChatMessage, error) {
	return instance.GetChatMessageSeqAfter(uid, seqAfter)
}
func GetChatMessageSeqSpan(uid int64, seq int64) (int, error) {
	return instance.GetChatMessageSeqSpan(uid, seq)
}
func AddOfflineMessage(uid int64, mid int64) error {
	return instance.AddOfflineMessage(uid, mid)
}
func GetOfflineMessage(uid int64) ([]*OfflineMessage, error) {
	return instance.GetOfflineMessage(uid)
}
func DelOfflineMessage(uid int64, mid []int64) error {
	return instance.DelOfflineMessage(uid, mid)
}
