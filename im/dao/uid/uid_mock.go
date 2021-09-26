package uid

import "time"

func Mock() {
	instance = &_gen{
		uid:     userIdStart + 1,
		sysUid:  systemIdStart + 1,
		tempUid: tempIdStart + 1,
	}
}

type _gen struct {
	uid     int64
	sysUid  int64
	tempUid int64
}

func (m *_gen) GenSysUid() int64 {
	m.sysUid += int64(time.Now().Second())
	m.sysUid++
	return m.sysUid
}

func (m *_gen) GenUid() int64 {
	m.uid += int64(time.Now().Second())
	m.uid++
	return m.uid
}

func (m *_gen) GenTempUid() int64 {
	m.tempUid += int64(time.Now().Second())
	m.tempUid++
	return m.tempUid
}
