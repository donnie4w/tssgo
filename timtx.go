// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"github.com/donnie4w/gofer/buffer"
	"github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/tssgo/stub"
)

type tx struct {
	account  string
	password string
	domain   string
}

func (ts *tx) ping() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMPING))
	return buf.Bytes()
}

func (ts *tx) login() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMAUTH))
	ab := stub.NewAuthBean()
	ab.Domain = &ts.domain
	ab.Username = &ts.account
	ab.Password = &ts.password
	buf.Write(thrift.TEncode(ab))
	return buf.Bytes()
}

func (ts *tx) subonline() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMSUB))
	ab := stub.NewAdmSubBean()
	ab.SubType = &SUB_ONLINE
	buf.Write(thrift.TEncode(ab))
	return buf.Bytes()
}
