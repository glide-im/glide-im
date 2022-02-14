package appdao

import (
	"go_im/pkg/db"
	"testing"
	"time"
)

func init() {
	db.Init()
}

func TestImpl_UpdateReleaseInfo(t *testing.T) {
	err := Impl.UpdateReleaseInfo(&ReleaseInfo{
		VersionCode: 1,
		VersionName: "1",
		UpdateAt:    time.Now().Unix(),
		DownloadUrl: "https://",
		Description: "none",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestImpl_GetReleaseInfo(t *testing.T) {
	info, err := Impl.GetReleaseInfo()
	if err != nil {
		t.Error(err)
	}
	t.Log(info)
}
