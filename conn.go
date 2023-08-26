package doh

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type Conn struct {
	io.Writer
	response func() (io.Reader, error)
	Reset    func()
}

var _ net.Conn = (*Conn)(nil)

func NewConn(client *http.Client, ctx context.Context, server string) (conn *Conn) {
	if client == nil {
		client = http.DefaultClient
	}
	body := new(bytes.Buffer)
	conn = &Conn{Writer: body}
	request := func() (reader io.Reader, err error) {
		link := fmt.Sprintf("https://%s/dns-query", server)
		body.Next(2) // skip uint16 length [2]byte
		req, err := http.NewRequest(http.MethodPost, link, body)
		if err != nil {
			return
		}
		if ctx != nil {
			req = req.WithContext(ctx)
		}
		req.Header.Set("Content-Type", "application/dns-message")
		req.Header.Add("Accept", "application/dns-message")
		resp, err := client.Do(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("expect code 200, but got %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		length := new(bytes.Buffer)
		binary.Write(length, binary.BigEndian, uint16(len(body)))
		reader = io.MultiReader(length, bytes.NewBuffer(body))
		return reader, nil
	}
	conn.Reset = func() {
		body.Reset()
		conn.response = sync.OnceValues(request)
	}
	conn.Reset()
	return conn
}

func (conn *Conn) Read(b []byte) (n int, err error) {
	reader, err := conn.response()
	if err != nil {
		return
	}
	n, err = reader.Read(b)
	return
}

func (conn *Conn) Close() error                  { return nil }
func (*Conn) LocalAddr() net.Addr                { return nil }
func (*Conn) RemoteAddr() net.Addr               { return nil }
func (*Conn) SetDeadline(t time.Time) error      { return nil }
func (*Conn) SetReadDeadline(t time.Time) error  { return nil }
func (*Conn) SetWriteDeadline(t time.Time) error { return nil }
