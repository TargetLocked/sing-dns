package dns

import (
	"github.com/sagernet/sing/common/buf"

	"github.com/miekg/dns"
)

func TruncateDNSMessage(request *dns.Msg, response *dns.Msg, frontHeadroom int) (*buf.Buffer, error) {
	maxLen := 512
	if edns0Option := request.IsEdns0(); edns0Option != nil {
		if udpSize := int(edns0Option.UDPSize()); udpSize > 512 {
			maxLen = udpSize
		}
	}
	response.Truncate(maxLen)
	buffer := buf.NewSize(frontHeadroom + 1 + maxLen)
	buffer.Advance(frontHeadroom)
	rawMessage, err := response.PackBuffer(buffer.FreeBytes())
	if err != nil {
		buffer.Release()
		return nil, err
	}
	buffer.Truncate(len(rawMessage))
	return buffer, nil
}