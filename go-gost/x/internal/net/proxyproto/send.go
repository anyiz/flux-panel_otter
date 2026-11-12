package proxyproto

import (
	"fmt"
	pp "github.com/pires/go-proxyproto"
	"net"
	"time"
)

// WriteClientHeader writes the header as tunnel payload when c is a relay
// connector. The exit forwards these bytes unchanged to the final target.
func WriteClientHeader(version int, src, dst net.Addr, c net.Conn) error {
	if version != 1 && version != 2 {
		return fmt.Errorf("invalid PROXY version")
	}
	header := pp.HeaderProxyFromAddrs(byte(version), src, dst)
	if err := c.SetWriteDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return err
	}
	defer c.SetWriteDeadline(time.Time{})
	_, err := header.WriteTo(c)
	return err
}
