package dao

const (
	TempIdStart = 1_000_000_000_000
	TempIdEnd   = TempIdStart + 1_000_000_000
)

func IsTempId(uid int64) bool {
	return TempIdStart < uid && uid < TempIdEnd
}

func GenUid() int64 {
	return 0
}

type TempIdGen struct {
}

func (g *TempIdGen) Obtain() int64 {
	return 0
}

func (g *TempIdGen) Recycle(uid int64) {

}
