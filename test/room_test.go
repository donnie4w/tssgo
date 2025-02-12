package test

import "testing"

func Test_CreateRoom(t *testing.T) {
	client := newClient()
	ack, _ := client.CreateRoom("UHuS8PoK2M9", "tlnet.top", "testroom", 1)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_RoomUser(t *testing.T) {
	client := newClient()
	r, _ := client.RoomUsers("KBn7xAvqsqR", "tlnet.top")
	t.Log(r)
}

func Test_AddRoom(t *testing.T) {
	client := newClient()
	r, _ := client.AddRoom("QdH6CCms5Ex", "tlnet.top", "KBn7xAvqsqR", "")
	t.Log(r)
}

func Test_PullInRoom(t *testing.T) {
	client := newClient()
	ack, _ := client.PullInRoom("UHuS8PoK2M9", "tlnet.top", "KBn7xAvqsqR", "QdH6CCms5Ex")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_KickRoom(t *testing.T) {
	client := newClient()
	ack, _ := client.KickRoom("UHuS8PoK2M9", "tlnet.top", "KBn7xAvqsqR", "QdH6CCms5Ex")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}
