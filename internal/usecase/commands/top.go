package commands

import (
	"fmt"
	"proxy-server-with-tg-admin/internal/helper"
	"runtime"
)

type top struct {
}

func (c *top) Id() string {
	return "top"
}

func (c *top) Arguments() []string {
	return []string{}
}

func (c *top) Run(args ...string) (string, error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	res := ""

	res += fmt.Sprintf("NumGoroutine = %v\n", runtime.NumGoroutine())
	res += fmt.Sprintf("Alloc = %s\n", helper.BytesFormat(m.Alloc))
	res += fmt.Sprintf("TotalAlloc = %s\n", helper.BytesFormat(m.TotalAlloc))
	res += fmt.Sprintf("Sys = %s\n", helper.BytesFormat(m.Sys))
	res += fmt.Sprintf("NumGC = %v\n", m.NumGC)

	return res, nil
}
