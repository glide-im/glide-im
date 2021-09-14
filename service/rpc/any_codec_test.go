package rpc

import (
	"github.com/stretchr/testify/assert"
	"go_im/service/pb"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestPB4AnyWrapperCodec_Decode(t *testing.T) {

	codec := PB4AnyWrapperCodec{}
	r := &pb.UidRequest{Uid: 1}
	b, err := codec.Encode(r)
	assert.Nil(t, err)

	r2 := &pb.UidRequest{}
	err = codec.Decode(b, r2)
	assert.Nil(t, err)
	assert.Equal(t, r.Uid, r2.Uid)
}

func TestPB4AnyWrapperCodec_Encode(t *testing.T) {

	codec := PB4AnyWrapperCodec{}
	r := &pb.UidRequest{Uid: 1}
	b, err := codec.Encode(r)
	assert.Nil(t, err)

	r2 := &pb.UidRequest{}
	err = proto.Unmarshal(b, r2)
	assert.Nil(t, err)
	assert.Equal(t, r2.Uid, r.Uid)
}
