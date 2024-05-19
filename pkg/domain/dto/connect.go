package dto

import (
	"connectrpc.com/connect"
)

type ConnectSpec struct {
	Procedure        string
	Description      string
	IdempotencyLevel connect.IdempotencyLevel
}

type ConnectSpecs []*ConnectSpec
