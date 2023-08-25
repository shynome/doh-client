## Intro

doh conn for [github.com/miekg/dns](https://github.com/miekg/dns)

```go
package main

import (
	"fmt"

	"github.com/miekg/dns"
	"github.com/shynome/doh-client"
)

func main() {
	q := dns.Question{
		Name:   dns.Fqdn("remoon.net"),
		Qtype:  dns.TypeA,
		Qclass: dns.ClassINET,
	}
	m := &dns.Msg{
		MsgHdr:   dns.MsgHdr{Id: dns.Id(), Opcode: dns.OpcodeQuery, RecursionDesired: true},
		Question: []dns.Question{q},
	}
	co := &dns.Conn{Conn: doh.NewConn(nil, nil, "1.1.1.1")}
	if err := co.WriteMsg(m); err != nil {
		panic(err)
	}
	m, err := co.ReadMsg()
	if err != nil {
		panic(err)
	}
	if len(m.Answer) == 0 {
		panic("answer length must greater than 0")
	}
	fmt.Println(m.Answer)
}

```
