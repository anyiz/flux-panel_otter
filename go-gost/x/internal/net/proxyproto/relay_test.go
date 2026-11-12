package proxyproto_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/go-gost/core/connector"
	corechain "github.com/go-gost/core/chain"
	"github.com/go-gost/core/handler"
	"github.com/go-gost/x/chain"
	relayconnector "github.com/go-gost/x/connector/relay"
	relayhandler "github.com/go-gost/x/handler/relay"
	"github.com/go-gost/x/internal/net/proxyproto"
	"github.com/go-gost/x/logger"
	"github.com/go-gost/x/metadata"
	pp "github.com/pires/go-proxyproto"
)

// Exercise production relay framing, including buffered and immediate handshakes.
// The landing server speaks first after the PROXY header, before client payload.
func TestHeaderThroughRealRelay(t *testing.T) {
	for _, noDelay := range []bool{false, true} {
		for _, version := range []int{1, 2} {
			t.Run(fmt.Sprintf("v%d/nodelay=%v", version, noDelay), func(t *testing.T) {
				target, e := net.Listen("tcp", "127.0.0.1:0")
				if e != nil {
					t.Fatal(e)
				}
				defer target.Close()
				done := make(chan error, 1)
				go func() {
					c, e := target.Accept()
					if e != nil {
						done <- e
						return
					}
					defer c.Close()
					c.SetDeadline(time.Now().Add(3 * time.Second))
					reader := bufio.NewReader(c)
					h, e := pp.Read(reader)
					if e != nil {
						done <- e
						return
					}
					if h.SourceAddr.String() != "192.0.2.7:1234" {
						done <- fmt.Errorf("wrong source: %v", h.SourceAddr)
						return
					}
					if _, e = c.Write([]byte("hello")); e != nil {
						done <- e
						return
					}
					b := make([]byte, 4)
					_, e = io.ReadFull(reader, b)
					if e == nil && string(b) != "ping" {
						e = fmt.Errorf("payload changed: %q", b)
					}
					done <- e
				}()
				log := logger.NewLogger(logger.OutputOption(io.Discard))
				md := metadata.NewMetadata(map[string]any{"nodelay": noDelay})
				h := relayhandler.NewHandler(handler.LoggerOption(log), handler.RouterOption(chain.NewRouter(corechain.LoggerRouterOption(log))))
				if e = h.Init(md); e != nil {
					t.Fatal(e)
				}
				if closer, ok := h.(io.Closer); ok {
					defer closer.Close()
				}
				relay, e := net.Listen("tcp", "127.0.0.1:0")
				if e != nil {
					t.Fatal(e)
				}
				defer relay.Close()
				go func() {
					c, e := relay.Accept()
					if e == nil {
						defer c.Close()
						h.Handle(context.Background(), c)
					}
				}()
				raw, e := net.Dial("tcp", relay.Addr().String())
				if e != nil {
					t.Fatal(e)
				}
				defer raw.Close()
				raw.SetDeadline(time.Now().Add(3 * time.Second))
				co := relayconnector.NewConnector(connector.LoggerOption(log))
				if e = co.Init(md); e != nil {
					t.Fatal(e)
				}
				stream, e := co.Connect(context.Background(), raw, "tcp", target.Addr().String())
				if e != nil {
					t.Fatal(e)
				}
				if e = proxyproto.WriteClientHeader(version, &net.TCPAddr{IP: net.ParseIP("192.0.2.7"), Port: 1234}, &net.TCPAddr{IP: net.ParseIP("203.0.113.7"), Port: 443}, stream); e != nil {
					t.Fatal(e)
				}
				greeting := make([]byte, 5)
				if _, e = io.ReadFull(stream, greeting); e != nil {
					t.Fatal(e)
				}
				if string(greeting) != "hello" {
					t.Fatal(string(greeting))
				}
				if _, e = stream.Write([]byte("ping")); e != nil {
					t.Fatal(e)
				}
				if e = <-done; e != nil {
					t.Fatal(e)
				}
			})
		}
	}
}
