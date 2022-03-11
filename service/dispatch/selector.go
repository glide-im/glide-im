package dispatch

import "context"

type dispatchSelector struct {
}

func (u *dispatchSelector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {

	return ""
}

func (u *dispatchSelector) UpdateServer(servers map[string]string) {

}
