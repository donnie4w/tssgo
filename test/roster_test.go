package test

import "testing"

func Test_Addroster(t *testing.T) {
	client := newClient()
	ack, _ := client.Addroster("QdH6CCms5Ex", "tlnet.top", "UHuS8PoK2M9", "hello!")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Blockroster(t *testing.T) {
	client := newClient()
	ack, _ := client.Blockroster("UHuS8PoK2M9", "tlnet.top", "QdH6CCms5Ex")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Rmroster(t *testing.T) {
	client := newClient()
	ack, _ := client.Rmroster("UHuS8PoK2M9", "tlnet.top", "QdH6CCms5Ex")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Roster(t *testing.T) {
	client := newClient()
	if r, err := client.Roster("UHuS8PoK2M9", "tlnet.top"); err == nil {
		t.Log(r.GetNodelist())
	}
}
