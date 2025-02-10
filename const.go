// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

type TIMTYPE byte

const (
	ADMPING         TIMTYPE = 11
	ADMAUTH         TIMTYPE = 12
	ADMSUB          TIMTYPE = 13
	BIGBINARYSTREAM TIMTYPE = 17
)

const (
	SOURCE_OS   int8 = 1
	SOURCE_USER int8 = 2
	SOURCE_ROOM int8 = 3
)

var (
	ONLINE   int8 = 1
	OFFLIINE int8 = 2
)

const (
	ORDER_INOF      int8 = 1
	ORDER_REVOKE    int8 = 2
	ORDER_BURN      int8 = 3
	ORDER_BUSINESS  int8 = 4
	ORDER_STREAM    int8 = 5
	ORDER_BIGSTRING int8 = 6
	ORDER_BIGBINARY int8 = 7
	ORDER_RESERVED  int8 = 30
)

var (
	SUB_ONLINE int8 = 1
)

var (
	SEP_BIN = byte(131)
)
