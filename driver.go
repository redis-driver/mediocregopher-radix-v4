package radix4

import (
	"context"
	"github.com/mediocregopher/radix/v4"
	"errors"
	redScript "github.com/redis-driver/script"
	"log"
	"strings"
)

type Client struct {
	Core radix.Client
	debug bool
}
// 为了防止调试之后忘记去掉 DebugOnce, 函数签名返回 error 可以让代码在编译期检查提示有错误未处理，实际上DebugOnce 永远返回 nil。
func (c *Client) DebugOnce() error {
	c.debug = true
	return nil
}
func (c *Client) logDebug(cmd []string) {
	if c.debug {
		c.debug = false
		log.Print("goclub/redis:(debug) exec: ", strings.Join(cmd, " "))
	}
}
func (c Client) RedisCommand(ctx context.Context, valuePtr interface{}, argv []string) (result struct { IsNil bool }, err error){
	c.logDebug(argv)
	data := radix.Maybe{Rcv: valuePtr}
	var moreArg []string
	if len(argv) >1 { moreArg = argv[1:] }
	err = c.Core.Do(ctx, radix.Cmd(&data, argv[0], moreArg...)) ; if err != nil {
		return
	}
	result.IsNil = data.Null
	return
}

func (c Client)  RedisScript (ctx context.Context, script redScript.Script) (err error){
	err = c.Core. Do(ctx, radix.NewEvalScript(script.Script).Cmd(script.ValuePtr, script.Keys, script.Argv...)) ; if err != nil {
		return
	}
	return
}


func (c Client)  Close () error {
	if c.Core == nil {
		return errors.New("radix client is nil can not close")
	}
	return c.Core.Close()
}

type StreamEntryFields [][2]string
func (data StreamEntryFields) Field(name string) (value string, hasValue bool) {
	for _, item := range data {
		if item[0] == name {
			return item[1], true
		}
	}
	return "", false
}
