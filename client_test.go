// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"github.com/donnie4w/tssgo/stub"
	"testing"
)

func testNewConnect() *Client {
	client, err := NewClient(false, "127.0.0.1:9002", "admin", "123", "tlnet.top")
	if err != nil {
		panic(err)
	}
	return client
}

// tom1: 32KUvu5Xqj5
// tom2: 6gwprA48tag
// tom3: ANpPiFdPw6A
func Test_Register(t *testing.T) {
	client := testNewConnect()
	ab := stub.NewAuthBean()
	name, pwd, domain := "tom3", "123", "tlnet.top"
	ab.Username, ab.Password, ab.Domain = &name, &pwd, &domain
	ack, _ := client.Register(ab)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_Token(t *testing.T) {
	client := testNewConnect()
	at := stub.NewAdmToken()
	name, pwd, domain := "tom1", "123", "tlnet.top"
	at.Name, at.Password, at.Domain = &name, &pwd, &domain
	ack, _ := client.Token(at)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

func Test_CreateRoom(t *testing.T) {
	client := testNewConnect()
	ack, _ := client.CreateRoom("32KUvu5Xqj5", "tlnet.top", "testroom", 1)
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}

// b5DMLPxo2RW
func Test_RoomUser(t *testing.T) {
	client := testNewConnect()
	r, _ := client.RoomUsers("b5DMLPxo2RW", "tlnet.top")
	t.Log(r)
}

func Test_AddRoom(t *testing.T) {
	client := testNewConnect()
	r, _ := client.AddRoom("6gwprA48tag", "tlnet.top", "b5DMLPxo2RW", "")
	t.Log(r)
}

func Test_PullInRoom(t *testing.T) {
	client := testNewConnect()
	ack, _ := client.PullInRoom("32KUvu5Xqj5", "tlnet.top", "b5DMLPxo2RW", "6gwprA48tag")
	t.Log(ack.GetOk(), ",", ack.GetT(), ",", ack.GetN(), ",", ack.GetTimType(), ",", ack.GetErrcode())
}
