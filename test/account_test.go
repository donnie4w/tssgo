package test

import (
	"github.com/donnie4w/tssgo/stub"
	"testing"
)

// tom1: 32KUvu5Xqj5
// tom2: 6gwprA48tag
// tom3: ANpPiFdPw6A
func Test_Register(t *testing.T) {
	client := newClient()
	ab := stub.NewAuthBean()
	name, pwd, domain := "tim1", "123", "tlnet.top"
	ab.Username, ab.Password, ab.Domain = &name, &pwd, &domain
	ack, _ := client.Register(ab)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Token(t *testing.T) {
	client := newClient()
	at := stub.NewAdmToken()
	name, pwd, domain := "tom1", "123", "tlnet.top"
	at.Name, at.Password, at.Domain = &name, &pwd, &domain
	ack, _ := client.Token(at)
	t.Log(ack.GetOk(), ",", ack.GetN(), ",", ack.GetN2(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Authroster(t *testing.T) {
	client := newClient()
	ack, _ := client.Authroster("QdH6CCms5Ex", "tlnet.top", "UHuS8PoK2M9")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Authgroupuser(t *testing.T) {
	client := newClient()
	ack, _ := client.Authgroupuser("QdH6CCms5Ex", "tlnet.top", "KBn7xAvqsqR")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_SysBlockUser(t *testing.T) {
	client := newClient()
	abu := stub.NewAdmSysBlockUser()
	abu.Nodelist = []string{"tim1"}
	bt := int64(60)
	abu.Blocktime = &bt
	ack, _ := client.SysBlockUser(abu)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}
