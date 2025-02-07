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

type tsstx struct {
	account  string
	password string
	domain   string
}

func (tt *tsstx) ping() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMPING))
	return buf.Bytes()
}

func (tt *tsstx) login() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMAUTH))
	ab := stub.NewAuthBean()
	ab.Domain = &tt.domain
	ab.Username = &tt.account
	ab.Password = &tt.password
	buf.Write(thrift.TEncode(ab))
	return buf.Bytes()
}

func (tt *tsstx) subonline() []byte {
	buf := buffer.NewBuffer()
	buf.WriteByte(byte(ADMSUB))
	ab := stub.NewAdmSubBean()
	ab.SubType = &SUB_ONLINE
	buf.Write(thrift.TEncode(ab))
	return buf.Bytes()
}

func (tt *tsstx) bigbinaryStream(vnode, fnode string, body []byte) []byte {
	buf := buffer.NewBufferWithCapacity(3 + len(vnode) + len(fnode) + len(body))
	buf.WriteByte(byte(BIGBINARYSTREAM))
	buf.WriteString(vnode)
	buf.WriteByte(SEP_BIN)
	buf.WriteString(fnode)
	buf.WriteByte(SEP_BIN)
	buf.Write(body)
	return buf.Bytes()
}
