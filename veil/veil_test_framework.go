package veil

import (
	"testing"
)

type Control interface {
	StartTest(t testing.TB)
	StopTest(t testing.TB)
}

type control struct {
}

func (c *control) StartTest(t testing.TB) {

}
func (c *control) StopTest(t testing.TB) {

}

func InitTestFramework(connFactory ConnectionFactory, serverFactory ServerFactory) Control {
	VeilInitClient(connFactory)
	VeilInitServer(serverFactory)

	return &control{}
}
