package main

import (
	kwekker_protobufs "github.com/googolplex-s6/kwekker-protobufs/kwek"
	"go.uber.org/zap"
	"kwekker-worker/util"
)

func Initialize(logger *zap.Logger, config *util.Config) {
	_ = make(chan *kwekker_protobufs.Kwek)
}
