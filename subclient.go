// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

import (
	"fmt"
	"github.com/donnie4w/gofer/thrift"
	"github.com/donnie4w/gofer/websocket"
	"github.com/donnie4w/tssgo/stub"
	"time"
)

type SubClient struct {
	addr         string
	pingCount    int
	handler      *websocket.Handler
	cfg          *websocket.Config
	ts           *tx
	isclose      bool
	SubHandler   func(bean *stub.AdmSubBean)
	OnLoginEvent func()
}

func NewSubClient(tls bool, addr string, account, password, domain string) (sc *SubClient) {
	if addr := formatUrl(addr, tls); addr != "" {
		sc = &SubClient{addr: addr, ts: &tx{account: account, password: password, domain: domain}}
		sc.defaultInit()
	}
	return
}

func (sc *SubClient) login() (err error) {
	if sc.handler != nil {
		err = sc.send(sc.ts.login())
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

func (sc *SubClient) init(cfg *websocket.Config) {
	cfg.OnError = func(_ *websocket.Handler, err error) {
		logger.Error("OnError:", err)
		sc.close()
		<-time.After(time.Second << 2)
		if !sc.isclose {
			sc.connect()
		}
	}
	cfg.OnMessage = func(_ *websocket.Handler, msg []byte) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(err)
			}
		}()
		sc.pingCount = 0
		sc.doMsg(TIMTYPE(msg[0]&0x7f), msg[1:])
	}
	sc.cfg = cfg
}

func (sc *SubClient) defaultInit() {
	sc.init(&websocket.Config{TimeOut: 10 * time.Second, Url: sc.addr + "/tim", Origin: "https://github.com/donnie4w/tim"})
}

func (sc *SubClient) SetTimeOut(t time.Duration) {
	sc.cfg.TimeOut = t
}

func (sc *SubClient) Connect() error {
	return sc.connect()
}

func (sc *SubClient) connect() (err error) {
	sc.pingCount = 0
	if sc.handler, err = websocket.NewHandler(sc.cfg); err == nil {
		err = sc.login()
		go sc.ping()
	} else if !sc.isclose {
		<-time.After(3 * time.Second)
		logger.Warn("reconnect")
		sc.connect()
	}
	<-time.After(time.Second)
	return
}

func (sc *SubClient) Close() (err error) {
	sc.isclose = true
	return sc.close()
}

func (sc *SubClient) send(bs []byte) (err error) {
	return sc.handler.Send(bs)
}

func (sc *SubClient) close() (err error) {
	if sc.handler != nil {
		err = sc.handler.Close()
	}
	return
}

func (sc *SubClient) doMsg(t TIMTYPE, bs []byte) {
	switch t {
	case ADMPING:
		if sc.pingCount > 0 {
			sc.pingCount--
		}
	case ADMAUTH:
		if ta, err := thrift.TDecode(bs, stub.NewAdmAck()); ta != nil && err == nil && ta.GetOk() {
			logger.Info("auth successful")
			if sc.OnLoginEvent != nil {
				sc.OnLoginEvent()
			}
		} else {
			logger.Error("auth failed:", sc.ts.account)
		}
	case ADMSUB:
		if sc.SubHandler != nil {
			if tm, _ := thrift.TDecode(bs, stub.NewAdmSubBean()); tm != nil {
				sc.SubHandler(tm)
			}
		}
	default:
		logger.Warn("undisposed >>>>>", t, " ,data length:", len(bs))
	}
}

func (sc *SubClient) ping() {
	defer recoverable(nil)
	ticker := time.NewTicker(15 * time.Second)
	for !sc.isclose {
		select {
		case <-ticker.C:
			sc.pingCount++
			if sc.isclose {
				goto END
			}
			if err := sc.send(sc.ts.ping()); err != nil || sc.pingCount > 3 {
				logger.Error("ping over count>>", sc.pingCount, err)
				sc.close()
				goto END
			}
		}
	}
END:
}

func (sc *SubClient) SubOnline() error {
	return sc.send(sc.ts.subonline())
}

func recoverable(err *error) {
	if e := recover(); e != nil {
		if err != nil {
			*err = fmt.Errorf("panic: %v", e)
		}
	}
}
