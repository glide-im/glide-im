package appdao

import (
	"github.com/glide-im/glideim/pkg/db"
	"testing"
	"time"
)

func init() {
	db.Init()
}

func TestImpl_UpdateReleaseInfo(t *testing.T) {
	err := Impl.UpdateReleaseInfo(&ReleaseInfo{
		VersionCode: 142,
		VersionName: "142",
		UpdateAt:    time.Now().Unix(),
		DownloadUrl: "https://github.91chi.fun//https://github.com//Glide-IM/Glide-IM-Android/releases/download/v1.4.1/glide_im_release_v1.4.1.apk",
		Description: "你没有体验过的全新版本\n快来下载体验吧!",
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
