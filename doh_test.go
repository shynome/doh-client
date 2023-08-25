package doh

import (
	"testing"

	"github.com/miekg/dns"
)

func TestDoH(t *testing.T) {
	q := dns.Question{
		Name:   dns.Fqdn("remoon.net"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}
	co := &dns.Conn{Conn: NewConn(nil, nil, "1.1.1.1")}
	if err := co.WriteMsg(m); err != nil {
		t.Error(err)
		return
	}
	m, err := co.ReadMsg()
	if err != nil {
		t.Error(err)
		return
	}
	if len(m.Answer) == 0 {
		t.Error("answer length must greater than 0")
		t.Error(m)
		return
	}
	t.Log(m.Answer)
}
