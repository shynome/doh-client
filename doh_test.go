package doh

import (
	"fmt"
	"testing"

	"github.com/miekg/dns"
)

func TestDoH(t *testing.T) {
	q := dns.Question{
		Name:   dns.Fqdn("remoon.net"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	for _, httpGet := range []bool{true, false} {
		conn := NewConn(nil, nil, "1.1.1.1")
		co := &dns.Conn{Conn: conn}
		conn.HttpGet = httpGet
		for i := 0; i < 2; i++ {
			err := func() (err error) {
				defer conn.Reset()
				m := &dns.Msg{
					MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
					Question: []dns.Question{q},
				}
				if err = co.WriteMsg(m); err != nil {
					return
				}
				m, err = co.ReadMsg()
				if err != nil {
					return
				}
				if len(m.Answer) == 0 {
					return fmt.Errorf("answer length must greater than 0")
				}
				t.Log(m.Id)
				t.Log(m.Answer)
				return
			}()
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
}
