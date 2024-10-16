// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"fmt"
	"github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/tssgo/stub"
	"testing"
	"time"
)

func Test_subclient(t *testing.T) {
	sc := NewSubClient(false, "127.0.0.1:9001", "admin", "123", "tlnet.top")
	sc.SubHandler = func(bean *stub.AdmSubBean) {
		fmt.Println(bean.GetSubType())
		switch bean.GetSubType() {
		case SUB_ONLINE:
			if asb, err := thrift.TDecode(bean.GetBs(), stub.NewAdmSubOnlineBean()); err == nil {
				fmt.Println("SUB_ONLINE>>>", asb.GetNode(), ":", asb.GetStatus())
			}
		}
	}
	sc.OnLoginEvent = func() {
		sc.SubOnline()
	}
	sc.Connect()
	<-time.After(500 * time.Second)
}
