package test

import (
	"github.com/donnie4w/tssgo"
	"github.com/donnie4w/tssgo/stub"
	"testing"
)

// 系统广播
func Test_TimMessageBroadcast(t *testing.T) {
	client := newClient()
	amb := stub.NewAdmMessageBroadcast()
	amb.Nodes = []string{"QdH6CCms5Ex", "UHuS8PoK2M9"}
	amb.Message = stub.NewAdmMessage()
	msg := "hello world 大家好"
	amb.Message.DataString = &msg
	ack, _ := client.TimMessageBroadcast(amb)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

// 系统广播状态
func Test_NewTimPresenceBroadcast(t *testing.T) {
	client := newClient()
	apb := stub.NewAdmPresenceBroadcast()
	apb.Nodes = []string{"QdH6CCms5Ex", "UHuS8PoK2M9"}
	apb.Presence = stub.NewAdmPresence()
	msg := "hello world 大家好"
	apb.Presence.Status = &msg
	ack, _ := client.TimPresenceBroadcast(apb)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_ProxyMessage(t *testing.T) {
	client := newClient()
	apm := stub.NewAdmProxyMessage()
	apm.Message = stub.NewAdmMessage()
	domain := "tlnet.top"
	from := "UHuS8PoK2M9"
	to := "QdH6CCms5Ex"
	msg := "hello world 大家好"
	apm.Message.MsType = tssgo.SOURCE_USER
	apm.Message.OdType = tssgo.ORDER_INOF
	apm.Message.Domain = &domain
	apm.Message.FromNode = &from
	apm.Message.ToNode = &to
	apm.Message.DataString = &msg
	if ack, _ := client.ProxyMessage(apm); ack != nil {
		t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
	}
}
