package appdao

import (
	"github.com/glide-im/glideim/pkg/db"
	"strconv"
)

const (
	keyReleaseInfo = "im:app:release"
)

var Impl AppDao = &impl{}

type AppDao interface {
	GetReleaseInfo() (*ReleaseInfo, error)
	UpdateReleaseInfo(r *ReleaseInfo) error
}

type impl struct {
}

func (i *impl) GetReleaseInfo() (*ReleaseInfo, error) {

	result, err := db.Redis.HGetAll(keyReleaseInfo).Result()
	if err != nil {
		return nil, err
	}
	info := &ReleaseInfo{
		VersionCode: 0,
		VersionName: result["version_name"],
		UpdateAt:    0,
		DownloadUrl: result["download_url"],
		Description: result["description"],
	}
	info.VersionCode, err = strconv.ParseInt(result["version_code"], 10, 64)
	if err != nil {
		return nil, err
	}
	info.UpdateAt, err = strconv.ParseInt(result["update_at"], 10, 64)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (i *impl) UpdateReleaseInfo(r *ReleaseInfo) error {
	set := db.Redis.HMSet(keyReleaseInfo, map[string]interface{}{
		"version_code": r.VersionCode,
		"version_name": r.VersionName,
		"update_at":    r.UpdateAt,
		"download_url": r.DownloadUrl,
		"description":  r.Description,
	})
	return set.Err()
}
