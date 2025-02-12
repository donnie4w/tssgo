// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/tssgo/stub"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var ConnectTimeout = 60 * time.Second
var SocketTimeout = 60 * time.Second
var transportFactory = thrift.NewTBufferedTransportFactory(1 << 13)
var logger = log.New(log.Writer(), "[tss]", log.LstdFlags)

type TssClient struct {
	admiface   *stub.AdmifaceClient
	transport  thrift.TTransport
	addr       string
	tls        bool
	authBean   *stub.AuthBean
	mux        *sync.Mutex
	pingCount  int32
	close      bool
	serverId   int64
	defaultCtx context.Context
}

func (tc *TssClient) connect(TLS bool, addr, name, pwd, domain string) (err error) {
	defer recoverable(&err)
	tc.tls, tc.addr, tc.authBean = TLS, addr, stub.NewAuthBean()
	tc.authBean.Username, tc.authBean.Password, tc.authBean.Domain = &name, &pwd, &domain
	tcf := &thrift.TConfiguration{ConnectTimeout: ConnectTimeout, SocketTimeout: SocketTimeout}
	var transport thrift.TTransport
	if tc.tls {
		tcf.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		transport = thrift.NewTSSLSocketConf(addr, tcf)
	} else {
		transport = thrift.NewTSocketConf(addr, tcf)
	}
	var useTransport thrift.TTransport
	if err == nil && transport != nil {
		if useTransport, err = transportFactory.GetTransport(transport); err == nil {
			if err = useTransport.Open(); err == nil {
				tc.addr = addr
				tc.transport = useTransport
				tc.admiface = stub.NewAdmifaceClientFactory(useTransport, thrift.NewTCompactProtocolFactoryConf(tcf))
				<-time.After(time.Second)
				tc.pingCount = 0
				err = tc.auth()
			}
		}
	}
	if err != nil {
		logger.Println("tss client to [", addr, "] Error:", err)
	}
	return
}

func (tc *TssClient) auth() error {
	if ack, err := tc.auth0(tc.authBean); err == nil {
		if ack.GetOk() {
			tc.serverId = ack.GetT()
			return nil
		}
	} else {
		return err
	}
	return errors.New("verification fail")
}

func (tc *TssClient) timer() {
	ticker := time.NewTicker(3 * time.Second)
	for !tc.close {
		select {
		case <-ticker.C:
			func() {
				defer recoverable(nil)
				if tc.pingCount > 5 {
					tc.reconnect()
					return
				}
				atomic.AddInt32(&tc.pingCount, 1)
				if ack, err := tc.ping(); err == nil && ack.GetOk() {
					atomic.AddInt32(&tc.pingCount, -1)
				}
			}()
		}
	}
}

func (tc *TssClient) ServerId() int64 {
	return tc.serverId
}

func (tc *TssClient) Close() (err error) {
	tc.close = true
	if tc.transport != nil {
		err = tc.transport.Close()
	}
	return
}

func (tc *TssClient) reconnect() error {
	logger.Println("tssgo reconnect")
	if tc.transport != nil {
		tc.transport.Close()
	}
	if !tc.close {
		return tc.connect(tc.tls, tc.addr, tc.authBean.GetUsername(), tc.authBean.GetPassword(), tc.authBean.GetDomain())
	} else {
		return errors.New("connection closed")
	}
}

func NewTssClient(tls bool, addr string, name, pwd, domain string) (client *TssClient, err error) {
	client = &TssClient{mux: &sync.Mutex{}, defaultCtx: context.TODO()}
	err = client.connect(tls, addr, name, pwd, domain)
	go client.timer()
	return
}

func (tc *TssClient) ping() (r *stub.AdmAck, err error) {
	defer recoverable(&err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Ping(tc.defaultCtx)
}

func (tc *TssClient) auth0(ab *stub.AuthBean) (r *stub.AdmAck, err error) {
	defer recoverable(&err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Auth(tc.defaultCtx, ab)
}

// ModifyPwd
// Parameters:
//   - Ab
func (tc *TssClient) ModifyPwd(node, oldpwd, newpwd, domain string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.ModifyPwd(tc.defaultCtx, node, oldpwd, newpwd, domain)
}

// Token
// Parameters:
//   - Atoken
func (tc *TssClient) Token(atoken *stub.AdmToken) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Token(tc.defaultCtx, atoken)
}

// TimMessageBroadcast
// Parameters:
//   - Amb
func (tc *TssClient) TimMessageBroadcast(amb *stub.AdmMessageBroadcast) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.TimMessageBroadcast(tc.defaultCtx, amb)
}

// TimPresenceBroadcast
// Parameters:
//   - Apb
func (tc *TssClient) TimPresenceBroadcast(apb *stub.AdmPresenceBroadcast) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.TimPresenceBroadcast(tc.defaultCtx, apb)
}

// ProxyMessage
// Parameters:
//   - Apm
func (tc *TssClient) ProxyMessage(apm *stub.AdmProxyMessage) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.ProxyMessage(tc.defaultCtx, apm)
}

// Register
// Parameters:
//   - Ab
func (tc *TssClient) Register(ab *stub.AuthBean) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Register(tc.defaultCtx, ab)
}

// ModifyUserInfo
// Parameters:
//   - Amui
func (tc *TssClient) ModifyUserInfo(amui *stub.AdmModifyUserInfo) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.ModifyUserInfo(tc.defaultCtx, amui)
}

// ModifyRoomInfo
// Parameters:
//   - Arb
func (tc *TssClient) ModifyRoomInfo(arb *stub.AdmRoomBean) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.ModifyRoomInfo(tc.defaultCtx, arb)
}

// SysBlockUser
// Parameters:
//   - Abu
func (tc *TssClient) SysBlockUser(abu *stub.AdmSysBlockUser) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.SysBlockUser(tc.defaultCtx, abu)
}

// OnlineUser
// Parameters:
//   - Au
func (tc *TssClient) OnlineUser(au *stub.AdmOnlineUser) (_r *stub.AdmTidList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.OnlineUser(tc.defaultCtx, au)
}

// Vroom
// Parameters:
//   - Avb
func (tc *TssClient) Vroom(avb *stub.AdmVroomBean) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Vroom(tc.defaultCtx, avb)
}

// Detect
// Parameters:
//   - Adb
func (tc *TssClient) Detect(adb *stub.AdmDetectBean) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Detect(tc.defaultCtx, adb)
}

// Roster
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) Roster(fromnode string, domain string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Roster(tc.defaultCtx, fromnode, domain)
}

// Addroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Msg
func (tc *TssClient) Addroster(fromnode string, domain string, tonode string, msg string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Addroster(tc.defaultCtx, fromnode, domain, tonode, msg)
}

// Rmroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
func (tc *TssClient) Rmroster(fromnode string, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Rmroster(tc.defaultCtx, fromnode, domain, tonode)
}

// Blockroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
func (tc *TssClient) Blockroster(fromnode string, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Blockroster(tc.defaultCtx, fromnode, domain, tonode)
}

// PullUserMessage
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Mid
//   - Limit
func (tc *TssClient) PullUserMessage(fromnode string, domain string, tonode string, mid, timeseries, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.PullUserMessage(tc.defaultCtx, fromnode, domain, tonode, mid, timeseries, limit)
}

// PullRoomMessage
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Mid
//   - Limit
func (tc *TssClient) PullRoomMessage(fromnode string, domain string, tonode string, mid, timeseries, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.PullRoomMessage(tc.defaultCtx, fromnode, domain, tonode, mid, timeseries, limit)
}

// OfflineMsg
// Parameters:
//   - Fromnode
//   - Domain
//   - Limit
func (tc *TssClient) OfflineMsg(fromnode string, domain string, limit int64) (_r *stub.AdmMessageList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.OfflineMsg(tc.defaultCtx, fromnode, domain, limit)
}

// DelOfflineMsg
//
//	Parameters:
//	 - Fromnode
//	 - Domain
//	 - ids
func (tc *TssClient) DelOfflineMsg(fromnode string, domain string, ids []int64) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.DelOfflineMsg(tc.defaultCtx, fromnode, domain, ids)
}

// UserRoom
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) UserRoom(fromnode string, domain string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.UserRoom(tc.defaultCtx, fromnode, domain)
}

// RoomUsers
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) RoomUsers(fromnode string, domain string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.RoomUsers(tc.defaultCtx, fromnode, domain)
}

// CreateRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - Topic : room topic
//   - Gtype : 1:GROUP_PRIVATE  2:GROUP_OPEN
func (tc *TssClient) CreateRoom(fromnode string, domain string, topic string, gtype int8) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.CreateRoom(tc.defaultCtx, fromnode, domain, topic, gtype)
}

// AddRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - Msg
func (tc *TssClient) AddRoom(fromnode string, domain string, roomNode string, msg string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.AddRoom(tc.defaultCtx, fromnode, domain, roomNode, msg)
}

// PullInRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (tc *TssClient) PullInRoom(fromnode string, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.PullInRoom(tc.defaultCtx, fromnode, domain, roomNode, toNode)
}

// RejectRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
//   - Msg
func (tc *TssClient) RejectRoom(fromnode string, domain string, roomNode string, toNode string, msg string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.RejectRoom(tc.defaultCtx, fromnode, domain, roomNode, toNode, msg)
}

// KickRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (tc *TssClient) KickRoom(fromnode string, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.KickRoom(tc.defaultCtx, fromnode, domain, roomNode, toNode)
}

// LeaveRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (tc *TssClient) LeaveRoom(fromnode string, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.LeaveRoom(tc.defaultCtx, fromnode, domain, roomNode)
}

// CancelRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (tc *TssClient) CancelRoom(fromnode string, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.CancelRoom(tc.defaultCtx, fromnode, domain, roomNode)
}

// BlockRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (tc *TssClient) BlockRoom(fromnode string, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.BlockRoom(tc.defaultCtx, fromnode, domain, roomNode)
}

// BlockRoomMember
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (tc *TssClient) BlockRoomMember(fromnode string, domain string, roomNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.BlockRoomMember(tc.defaultCtx, fromnode, domain, roomNode, toNode)
}

// BlockRosterList
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) BlockRosterList(fromnode string, domain string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.BlockRosterList(tc.defaultCtx, fromnode, domain)
}

// BlockRoomList
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) BlockRoomList(fromnode string, domain string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.BlockRoomList(tc.defaultCtx, fromnode, domain)
}

// BlockRoomMemberlist
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (tc *TssClient) BlockRoomMemberlist(fromnode string, domain string, roomNode string) (_r *stub.AdmNodeList, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.BlockRoomMemberlist(tc.defaultCtx, fromnode, domain, roomNode)
}

// VirtualroomRegister
// Parameters:
//   - Fromnode
//   - Domain
func (tc *TssClient) VirtualroomRegister(fromnode string, domain string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomRegister(tc.defaultCtx, fromnode, domain)
}

// VirtualroomRemove
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
func (tc *TssClient) VirtualroomRemove(fromnode string, domain string, vNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomRemove(tc.defaultCtx, fromnode, domain, vNode)
}

// VirtualroomAddAuth
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
//   - ToNode
func (tc *TssClient) VirtualroomAddAuth(fromnode string, domain string, vNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomAddAuth(tc.defaultCtx, fromnode, domain, vNode, toNode)
}

// VirtualroomDelAuth
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
//   - ToNode
func (tc *TssClient) VirtualroomDelAuth(fromnode string, domain string, vNode string, toNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomDelAuth(tc.defaultCtx, fromnode, domain, vNode, toNode)
}

// VirtualroomSub
// Parameters:
//   - Wsid
//   - Fromnode
//   - Domain
//   - VNode
//   - SubType
func (tc *TssClient) VirtualroomSub(wsid int64, fromnode string, domain string, vNode string, subType int8) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomSub(tc.defaultCtx, wsid, fromnode, domain, vNode, subType)
}

// VirtualroomUnSub
// Parameters:
//   - Wsid
//   - Fromnode
//   - Domain
//   - VNode
func (tc *TssClient) VirtualroomUnSub(wsid int64, fromnode string, domain string, vNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.VirtualroomUnSub(tc.defaultCtx, wsid, fromnode, domain, vNode)
}

func (tc *TssClient) Authroster(fromnode string, domain string, tonode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Authroster(tc.defaultCtx, fromnode, domain, tonode)
}

func (tc *TssClient) Authgroupuser(fromnode string, domain string, roomNode string) (_r *stub.AdmAck, _err error) {
	defer recoverable(&_err)
	tc.mux.Lock()
	defer tc.mux.Unlock()
	return tc.admiface.Authgroupuser(tc.defaultCtx, fromnode, domain, roomNode)
}
