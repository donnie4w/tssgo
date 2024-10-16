// Copyright (c) , donnie <donnie4w@gmail.com>
// All rights reserved.
//
// github.com/donnie4w/tssgo

package tssgo

type TIMTYPE byte

const (
	ADMPING TIMTYPE = 11
	ADMAUTH TIMTYPE = 12
	ADMSUB  TIMTYPE = 13
)

var (
	SUB_ONLINE int8 = 1
)
