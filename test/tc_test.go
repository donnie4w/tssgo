// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package test

import (
	"fmt"
	"github.com/donnie4w/tssgo"
)

func newClient() *tssgo.TssClient {
	client, err := tssgo.NewTssClient(false, "127.0.0.1:40000", "admin", "123", "tlnet.top")
	if err != nil {
		panic(err)
	}
	fmt.Println("tssgo connect success:", client.ServerId())
	return client
}
