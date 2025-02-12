// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package test

import (
	"github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/tssgo"
	"github.com/donnie4w/tssgo/stub"
	"testing"
	"time"
)

func Test_subclient(t *testing.T) {
	sc := tssgo.NewWsClient(false, "127.0.0.1:30000", "admin", "123", "tlnet.top")
	sc.SubHandler = func(bean *stub.AdmSubBean) {
		switch bean.GetSubType() {
		case tssgo.SUB_ONLINE:
			if asb, err := thrift.TDecode(bean.GetBs(), stub.NewAdmSubOnlineBean()); err == nil {
				t.Log("SUB_ONLINE>>>", asb.GetNode(), ":", asb.GetStatus())
			}
		default:
			t.Fatal("unknown sub type:", bean.GetSubType())
		}
	}
	sc.OnLoginEvent = func() {
		sc.SubOnline()
	}
	sc.BigbinaryStreamHandler = func(ack *stub.AdmAck) {
		if !ack.GetOk() {
			t.Logf("server Id:%d do not subscribe to this virtual room%s", ack.GetT(), ack.GetN())
		}
	}
	sc.Connect()
	time.Sleep(500 * time.Second)
}

func Test_BigbinaryStream(t *testing.T) {
	sc := tssgo.NewWsClient(false, "127.0.0.1:30001", "admin", "123", "tlnet.top")
	sc.Connect()
	<-time.After(3 * time.Second)
	bs := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sc.BigbinaryStream("5UuHvyshZND", "UHuS8PoK2M9", bs)
}
