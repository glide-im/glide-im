package broker

import (
	"go_im/pkg/rpc"
	"go_im/service/group_messaging"
)

type Client struct {
	*group_messaging.Client
}

func NewClient(options *rpc.ClientOptions) (*Client, error) {
	ret := &Client{}
	var err error

	// proxy of group messaging
	client, err := group_messaging.NewClient(options)
	if err != nil {
		return nil, err
	}
	ret.Client = client
	return ret, nil
}
