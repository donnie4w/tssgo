// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/donnie4w/simplelog/logging"
	. "github.com/donnie4w/tssgo/stub"
)

var ConnectTimeout = 60 * time.Second
var SocketTimeout = 60 * time.Second
var transportFactory = thrift.NewTBufferedTransportFactory(1 << 13)
var logger = logging.NewLogger().SetFormat(logging.FORMAT_LEVELFLAG | logging.FORMAT_TIME | logging.FORMAT_DATE)

type Client struct {
	admiface  *AdmifaceClient
	transport thrift.TTransport
	addr      string
	tls       bool
	authBean  *AuthBean
	mux       *sync.Mutex
	pingCount int32
	close     bool
}

func (cli *Client) connect(TLS bool, addr, name, pwd, domain string) (err error) {
	defer recoverable(&err)
	cli.tls, cli.addr, cli.authBean = TLS, addr, NewAuthBean()
	cli.authBean.Username, cli.authBean.Password, cli.authBean.Domain = &name, &pwd, &domain
	tcf := &thrift.TConfiguration{ConnectTimeout: ConnectTimeout, SocketTimeout: SocketTimeout}
	var transport thrift.TTransport
	if cli.tls {
		tcf.TLSConfig = &tls.Config{InsecureSkipVerify: true}
		transport = thrift.NewTSSLSocketConf(addr, tcf)
	} else {
		transport = thrift.NewTSocketConf(addr, tcf)
	}
	var useTransport thrift.TTransport
	if err == nil && transport != nil {
		if useTransport, err = transportFactory.GetTransport(transport); err == nil {
			if err = useTransport.Open(); err == nil {
				cli.addr = addr
				cli.transport = useTransport
				cli.admiface = NewAdmifaceClientFactory(useTransport, thrift.NewTCompactProtocolFactoryConf(tcf))
				<-time.After(time.Second)
				cli.pingCount = 0
				err = cli.auth()
			}
		}
	}
	if err != nil {
		logger.Error("tss client to [", addr, "] Error:", err)
	}
	return
}

func (cli *Client) auth() error {
	if ack, err := cli.auth0(cli.authBean); err == nil {
		if ack.GetOk() {
			return nil
		}
	} else {
		return err
	}
	return errors.New("verification fail")
}

func (cli *Client) timer() {
	ticker := time.NewTicker(3 * time.Second)
	for !cli.close {
		select {
		case <-ticker.C:
			func() {
				defer recoverable(nil)
				if cli.pingCount > 5 {
					cli.reconnect()
					return
				}
				atomic.AddInt32(&cli.pingCount, 1)
				if ack, err := cli.ping(); err == nil && ack.GetOk() {
					atomic.AddInt32(&cli.pingCount, -1)
				}
			}()
		}
	}
}

func (cli *Client) Close() (err error) {
	cli.close = true
	if cli.transport != nil {
		err = cli.transport.Close()
	}
	return
}

func (cli *Client) reconnect() error {
	logger.Warn("tssgo reconnect")
	if cli.transport != nil {
		cli.transport.Close()
	}
	if !cli.close {
		return cli.connect(cli.tls, cli.addr, cli.authBean.GetUsername(), cli.authBean.GetPassword(), cli.authBean.GetDomain())
	} else {
		return errors.New("connection closed")
	}
}

func NewClient(tls bool, addr string, name, pwd, domain string) (client *Client, err error) {
	client = &Client{mux: &sync.Mutex{}}
	err = client.connect(tls, addr, name, pwd, domain)
	go client.timer()
	return
}

func (cli *Client) ping() (r *AdmAck, err error) {
	defer recoverable(&err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Ping(context.TODO())
}

func (cli *Client) auth0(ab *AuthBean) (r *AdmAck, err error) {
	defer recoverable(&err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Auth(context.TODO(), ab)
}

// ModifyPwd
// Parameters:
//   - Ab
func (cli *Client) ModifyPwd(node, oldpwd, newpwd, domain string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.ModifyPwd(context.TODO(), node, oldpwd, newpwd, domain)
}

// Token
// Parameters:
//   - Atoken
func (cli *Client) Token(atoken *AdmToken) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Token(context.TODO(), atoken)
}

// TimMessageBroadcast
// Parameters:
//   - Amb
func (cli *Client) TimMessageBroadcast(amb *AdmMessageBroadcast) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.TimMessageBroadcast(context.TODO(), amb)
}

// TimPresenceBroadcast
// Parameters:
//   - Apb
func (cli *Client) TimPresenceBroadcast(apb *AdmPresenceBroadcast) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.TimPresenceBroadcast(context.TODO(), apb)
}

// ProxyMessage
// Parameters:
//   - Apm
func (cli *Client) ProxyMessage(apm *AdmProxyMessage) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.ProxyMessage(context.TODO(), apm)
}

// Register
// Parameters:
//   - Ab
func (cli *Client) Register(ab *AuthBean) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Register(context.TODO(), ab)
}

// ModifyUserInfo
// Parameters:
//   - Amui
func (cli *Client) ModifyUserInfo(amui *AdmModifyUserInfo) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.ModifyUserInfo(context.TODO(), amui)
}

// ModifyRoomInfo
// Parameters:
//   - Arb
func (cli *Client) ModifyRoomInfo(arb *AdmRoomBean) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.ModifyRoomInfo(context.TODO(), arb)
}

// SysBlockUser
// Parameters:
//   - Abu
func (cli *Client) SysBlockUser(abu *AdmSysBlockUser) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.SysBlockUser(context.TODO(), abu)
}

// SysBlockList
func (cli *Client) SysBlockList() (_r *AdmSysBlockList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.SysBlockList(context.TODO())
}

// OnlineUser
// Parameters:
//   - Au
func (cli *Client) OnlineUser(au *AdmOnlineUser) (_r *AdmTidList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.OnlineUser(context.TODO(), au)
}

// Vroom
// Parameters:
//   - Avb
func (cli *Client) Vroom(avb *AdmVroomBean) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Vroom(context.TODO(), avb)
}

// Detect
// Parameters:
//   - Adb
func (cli *Client) Detect(adb *AdmDetectBean) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Detect(context.TODO(), adb)
}

// Roster
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) Roster(fromnode string, domain string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Roster(context.TODO(), fromnode, domain)
}

// Addroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Msg
func (cli *Client) Addroster(fromnode string, domain string, tonode string, msg string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Addroster(context.TODO(), fromnode, domain, tonode, msg)
}

// Rmroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
func (cli *Client) Rmroster(fromnode string, domain string, tonode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Rmroster(context.TODO(), fromnode, domain, tonode)
}

// Blockroster
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
func (cli *Client) Blockroster(fromnode string, domain string, tonode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.Blockroster(context.TODO(), fromnode, domain, tonode)
}

// PullUserMessage
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Mid
//   - Limit
func (cli *Client) PullUserMessage(fromnode string, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.PullUserMessage(context.TODO(), fromnode, domain, tonode, mid, limit)
}

// PullRoomMessage
// Parameters:
//   - Fromnode
//   - Domain
//   - Tonode
//   - Mid
//   - Limit
func (cli *Client) PullRoomMessage(fromnode string, domain string, tonode string, mid int64, limit int64) (_r *AdmMessageList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.PullRoomMessage(context.TODO(), fromnode, domain, tonode, mid, limit)
}

// OfflineMsg
// Parameters:
//   - Fromnode
//   - Domain
//   - Limit
func (cli *Client) OfflineMsg(fromnode string, domain string, limit int64) (_r *AdmMessageList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.OfflineMsg(context.TODO(), fromnode, domain, limit)
}

// DelOfflineMsg
//
//	Parameters:
//	 - Fromnode
//	 - Domain
//	 - Ids
func (cli *Client) DelOfflineMsg(fromnode string, domain string, ids []int64) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.DelOfflineMsg(context.TODO(), fromnode, domain, ids)
}

// UserRoom
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) UserRoom(fromnode string, domain string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.UserRoom(context.TODO(), fromnode, domain)
}

// RoomUsers
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) RoomUsers(fromnode string, domain string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.RoomUsers(context.TODO(), fromnode, domain)
}

// CreateRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - Topic : room topic
//   - Gtype : 1:GROUP_PRIVATE  2:GROUP_OPEN
func (cli *Client) CreateRoom(fromnode string, domain string, topic string, gtype int8) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.CreateRoom(context.TODO(), fromnode, domain, topic, gtype)
}

// AddRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - Msg
func (cli *Client) AddRoom(fromnode string, domain string, roomNode string, msg string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.AddRoom(context.TODO(), fromnode, domain, roomNode, msg)
}

// PullInRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (cli *Client) PullInRoom(fromnode string, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.PullInRoom(context.TODO(), fromnode, domain, roomNode, toNode)
}

// RejectRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
//   - Msg
func (cli *Client) RejectRoom(fromnode string, domain string, roomNode string, toNode string, msg string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.RejectRoom(context.TODO(), fromnode, domain, roomNode, toNode, msg)
}

// KickRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (cli *Client) KickRoom(fromnode string, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.KickRoom(context.TODO(), fromnode, domain, roomNode, toNode)
}

// LeaveRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (cli *Client) LeaveRoom(fromnode string, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.LeaveRoom(context.TODO(), fromnode, domain, roomNode)
}

// CancelRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (cli *Client) CancelRoom(fromnode string, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.CancelRoom(context.TODO(), fromnode, domain, roomNode)
}

// BlockRoom
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (cli *Client) BlockRoom(fromnode string, domain string, roomNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.BlockRoom(context.TODO(), fromnode, domain, roomNode)
}

// BlockRoomMember
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
//   - ToNode
func (cli *Client) BlockRoomMember(fromnode string, domain string, roomNode string, toNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.BlockRoomMember(context.TODO(), fromnode, domain, roomNode, toNode)
}

// BlockRosterList
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) BlockRosterList(fromnode string, domain string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.BlockRosterList(context.TODO(), fromnode, domain)
}

// BlockRoomList
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) BlockRoomList(fromnode string, domain string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.BlockRoomList(context.TODO(), fromnode, domain)
}

// BlockRoomMemberlist
// Parameters:
//   - Fromnode
//   - Domain
//   - RoomNode
func (cli *Client) BlockRoomMemberlist(fromnode string, domain string, roomNode string) (_r *AdmNodeList, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.BlockRoomMemberlist(context.TODO(), fromnode, domain, roomNode)
}

// VirtualroomRegister
// Parameters:
//   - Fromnode
//   - Domain
func (cli *Client) VirtualroomRegister(fromnode string, domain string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomRegister(context.TODO(), fromnode, domain)
}

// VirtualroomRemove
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
func (cli *Client) VirtualroomRemove(fromnode string, domain string, vNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomRemove(context.TODO(), fromnode, domain, vNode)
}

// VirtualroomAddAuth
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
//   - ToNode
func (cli *Client) VirtualroomAddAuth(fromnode string, domain string, vNode string, toNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomAddAuth(context.TODO(), fromnode, domain, vNode, toNode)
}

// VirtualroomDelAuth
// Parameters:
//   - Fromnode
//   - Domain
//   - VNode
//   - ToNode
func (cli *Client) VirtualroomDelAuth(fromnode string, domain string, vNode string, toNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomDelAuth(context.TODO(), fromnode, domain, vNode, toNode)
}

// VirtualroomSub
// Parameters:
//   - Wsid
//   - Fromnode
//   - Domain
//   - VNode
//   - SubType
func (cli *Client) VirtualroomSub(wsid int64, fromnode string, domain string, vNode string, subType int8) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomSub(context.TODO(), wsid, fromnode, domain, vNode, subType)
}

// VirtualroomUnSub
// Parameters:
//   - Wsid
//   - Fromnode
//   - Domain
//   - VNode
func (cli *Client) VirtualroomUnSub(wsid int64, fromnode string, domain string, vNode string) (_r *AdmAck, _err error) {
	defer recoverable(&_err)
	cli.mux.Lock()
	defer cli.mux.Unlock()
	return cli.admiface.VirtualroomUnSub(context.TODO(), wsid, fromnode, domain, vNode)
}
