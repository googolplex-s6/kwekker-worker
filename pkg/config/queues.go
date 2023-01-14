package config

import (
	kwekproto "github.com/googolplex-s6/kwekker-protobufs/v3/kwek"
	userproto "github.com/googolplex-s6/kwekker-protobufs/v3/user"
	"google.golang.org/protobuf/proto"
)

type QueueData struct {
	Exchange string
	Type     proto.Message
}

type Queues map[string]QueueData

const kwekExchange = "kwek-exchange"
const userExchange = "user-exchange"

var QueueList = Queues{
	"kwek.create": {Exchange: kwekExchange, Type: &kwekproto.CreateKwek{}},
	"kwek.update": {Exchange: kwekExchange, Type: &kwekproto.UpdateKwek{}},
	"kwek.delete": {Exchange: kwekExchange, Type: &kwekproto.DeleteKwek{}},
	"user.create": {Exchange: userExchange, Type: &userproto.CreateUser{}},
	"user.update": {Exchange: userExchange, Type: &userproto.UpdateUser{}},
	"user.delete": {Exchange: userExchange, Type: &userproto.DeleteUser{}},
}
