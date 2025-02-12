// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"errors"
	"fmt"
	"github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/gofer/websocket"
	"github.com/donnie4w/tssgo/stub"
	"time"
)

type WsClient struct {
	addr                   string
	pingCount              int
	handler                *websocket.Handler
	cfg                    *websocket.Config
	ts                     *tsstx
	isclose                bool
	SubHandler             func(bean *stub.AdmSubBean)
	BigbinaryStreamHandler func(ack *stub.AdmAck)
	OnLoginEvent           func()
	err                    error
}

func NewWsClient(tls bool, addr string, account, password, domain string) (sc *WsClient) {
	if addr := formatUrl(addr, tls); addr != "" {
		sc = &WsClient{addr: addr, ts: &tsstx{account: account, password: password, domain: domain}}
		sc.defaultInit()
	}
	return
}

func (wc *WsClient) login() (err error) {
	if wc.handler != nil {
		err = wc.send(wc.ts.login())
	}
	return
}

func formatUrl(addr string, tls bool) (url string) {
	if tls {
		url = fmt.Sprint("wss://", addr)
	} else {
		url = fmt.Sprint("ws://", addr)
	}
	return
}

func (wc *WsClient) init(cfg *websocket.Config) {
	cfg.OnError = func(_ *websocket.Handler, err error) {
		logger.Println("OnError:", err)
		wc.close()
		<-time.After(time.Second << 2)
		if !wc.isclose {
			wc.connect()
		}
	}
	cfg.OnMessage = func(_ *websocket.Handler, msg []byte) {
		defer func() {
			if err := recover(); err != nil {
				logger.Println(err)
			}
		}()
		wc.pingCount = 0
		wc.doMsg(TIMTYPE(msg[0]&0x7f), msg[1:])
	}
	wc.cfg = cfg
}

func (wc *WsClient) defaultInit() {
	wc.init(&websocket.Config{TimeOut: 10 * time.Second, Url: wc.addr + "/tim", Origin: "https://github.com/donnie4w/tim"})
}

func (wc *WsClient) SetTimeOut(t time.Duration) {
	wc.cfg.TimeOut = t
}

func (wc *WsClient) Connect() error {
	return wc.connect()
}

func (wc *WsClient) Err() (err error) {
	return wc.err
}

func (wc *WsClient) connect() (err error) {
	wc.pingCount = 0
	if wc.handler, err = websocket.NewHandler(wc.cfg); err == nil {
		err = wc.login()
		go wc.ping()
	} else if !wc.isclose {
		<-time.After(3 * time.Second)
		logger.Println("reconnect")
		wc.connect()
	}
	<-time.After(time.Second)
	return
}

func (wc *WsClient) Close() (err error) {
	wc.isclose = true
	return wc.close()
}

func (wc *WsClient) send(bs []byte) (err error) {
	return wc.handler.Send(bs)
}

func (wc *WsClient) close() (err error) {
	defer recoverable(&err)
	if wc.handler != nil {
		err = wc.handler.Close()
	}
	return
}

func (wc *WsClient) doMsg(t TIMTYPE, bs []byte) {
	if wc.pingCount > 0 {
		wc.pingCount = 0
	}
	switch t {
	case ADMPING:
	case ADMAUTH:
		if ta, err := thrift.TDecode(bs, stub.NewAdmAck()); ta != nil && err == nil && ta.GetOk() {
			wc.err = nil
			logger.Println("Auth successful")
			if wc.OnLoginEvent != nil {
				wc.OnLoginEvent()
			}
		} else {
			s := fmt.Sprintf("Auth failed:", wc.ts.account, ",", wc.ts.password, ",", wc.ts.domain)
			wc.err = errors.New(s)
			logger.Println(s)
		}
	case ADMSUB:
		if wc.SubHandler != nil {
			if tm, _ := thrift.TDecode(bs, stub.NewAdmSubBean()); tm != nil {
				wc.SubHandler(tm)
			}
		}
	case BIGBINARYSTREAM:
		if wc.BigbinaryStreamHandler != nil {
			if ack, _ := thrift.TDecode(bs, stub.NewAdmAck()); ack != nil {
				wc.BigbinaryStreamHandler(ack)
			}
		}
	default:
		logger.Println("undisposed type:", t, " ,data length:", len(bs))
	}
}

func (wc *WsClient) ping() {
	defer recoverable(nil)
	ticker := time.NewTicker(15 * time.Second)
	for !wc.isclose {
		select {
		case <-ticker.C:
			wc.pingCount++
			if wc.isclose {
				goto END
			}
			if err := wc.send(wc.ts.ping()); err != nil || wc.pingCount > 3 {
				logger.Println("ping over count>>", wc.pingCount, err)
				wc.close()
				goto END
			}
		}
	}
END:
}

func (wc *WsClient) SubOnline() error {
	return wc.send(wc.ts.subonline())
}

func (wc *WsClient) BigbinaryStream(vnode, fnode string, body []byte) error {
	return wc.send(wc.ts.bigbinaryStream(vnode, fnode, body))
}

func recoverable(err *error) {
	if e := recover(); e != nil {
		if err != nil {
			*err = fmt.Errorf("panic: %v", e)
		}
	}
}
