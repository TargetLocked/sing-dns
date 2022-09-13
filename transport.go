package dns

import (
	"context"
	"net/netip"
	"net/url"

	E "github.com/sagernet/sing/common/exceptions"
	N "github.com/sagernet/sing/common/network"

	"github.com/miekg/dns"
)

type TransportConstructor = func(ctx context.Context, dialer N.Dialer, link string) (Transport, error)

type Transport interface {
	Start() error
	Close() error
	Raw() bool
	Exchange(ctx context.Context, message *dns.Msg) (*dns.Msg, error)
	Lookup(ctx context.Context, domain string, strategy DomainStrategy) ([]netip.Addr, error)
}

var transports map[string]TransportConstructor

func RegisterTransport(schemes []string, constructor TransportConstructor) {
	if transports == nil {
		transports = make(map[string]TransportConstructor)
	}
	for _, scheme := range schemes {
		transports[scheme] = constructor
	}
}

func CreateTransport(ctx context.Context, dialer N.Dialer, address string) (Transport, error) {
	constructor := transports[address]
	if constructor == nil {
		serverURL, err := url.Parse(address)
		if err == nil {
			constructor = transports[serverURL.Scheme]
		}
	}
	if constructor == nil {
		return nil, E.New("unknown DNS server format: " + address)
	}
	return constructor(ctx, dialer, address)
}
