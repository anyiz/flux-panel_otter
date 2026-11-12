package proxyproto

import (
	"bufio"
	pp "github.com/pires/go-proxyproto"
	"io"
	"net"
	"testing"
)

func TestHeaderPrecedesUnchangedPayload(t *testing.T) {
	for _, version := range []int{1, 2} {
		a, b := net.Pipe()
		done := make(chan error, 1)
		go func() {
			defer a.Close()
			e := WriteClientHeader(version, &net.TCPAddr{IP: net.ParseIP("192.0.2.7"), Port: 1234}, &net.TCPAddr{IP: net.ParseIP("203.0.113.2"), Port: 443}, a)
			if e == nil {
				_, e = a.Write([]byte("TLS-or-proxy-payload"))
			}
			done <- e
		}()
		reader := bufio.NewReader(b)
		h, e := pp.Read(reader)
		if e != nil {
			t.Fatal(e)
		}
		if h.SourceAddr.String() != "192.0.2.7:1234" {
			t.Fatal(h.SourceAddr)
		}
		data, e := io.ReadAll(reader)
		if e != nil {
			t.Fatal(e)
		}
		if string(data) != "TLS-or-proxy-payload" {
			t.Fatal(string(data))
		}
		if e = <-done; e != nil {
			t.Fatal(e)
		}
		b.Close()
	}
}
func TestFailedHeaderIsNotIgnored(t *testing.T) {
	a, b := net.Pipe()
	b.Close()
	defer a.Close()
	if e := WriteClientHeader(2, &net.TCPAddr{}, &net.TCPAddr{}, a); e == nil {
		t.Fatal("write failure ignored")
	}
}
