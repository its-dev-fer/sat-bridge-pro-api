package satws

import (
	"fmt"
	"strings"

	"github.com/InsaneTreset/nibus-sat-ws/credential"
	"github.com/InsaneTreset/nibus-sat-ws/service"
	"github.com/InsaneTreset/nibus-sat-ws/transport"
)

func NewCFDIClient(cer, key []byte, pass string) (*service.Client, error) {
	return NewClient(cer, key, pass, "cfdi", "")
}

func NewClient(cer, key []byte, pass, kind, proxyURL string) (*service.Client, error) {
	st, err := ParseServiceType(kind)
	if err != nil {
		return nil, err
	}
	cred, err := credential.NewFromBytes(cer, key, pass)
	if err != nil {
		return nil, err
	}
	if err := cred.Validate(); err != nil {
		return nil, err
	}
	var opts []service.Option
	if proxyURL != "" {
		hc, err := transport.New(proxyURL, 0)
		if err != nil {
			return nil, err
		}
		opts = append(opts, service.WithHTTPClient(hc))
	}
	return service.New(cred, st, opts...), nil
}

func ParseServiceType(kind string) (service.ServiceType, error) {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case "", "cfdi":
		return service.CFDI, nil
	case "retenciones":
		return service.Retenciones, nil
	default:
		return 0, fmt.Errorf("unknown SAT service %q", kind)
	}
}

func ParseDownloadType(v string) (service.DownloadType, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "received", "recibidas":
		return service.Received, nil
	case "issued", "emitidas":
		return service.Issued, nil
	default:
		return 0, fmt.Errorf("unknown download type %q", v)
	}
}

func ParseRequestType(v string) (service.RequestType, error) {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "cfdi":
		return service.RequestTypeCFDI, nil
	case "metadata":
		return service.RequestTypeMetadata, nil
	default:
		return 0, fmt.Errorf("unknown request type %q", v)
	}
}
