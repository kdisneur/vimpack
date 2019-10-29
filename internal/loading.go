package internal

import (
	"fmt"
)

type Loading int

const (
	LoadingStart Loading = iota
	LoadingOnDemand
)

func (l Loading) String() string {
	switch l {
	case LoadingStart:
		return "start"
	case LoadingOnDemand:
		return "opt"
	default:
		panic(fmt.Sprintf("impossible Loading value: %#+v", l))
	}
}
