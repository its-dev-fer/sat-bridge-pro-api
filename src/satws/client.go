package satws

import (
	"github.com/InsaneTreset/nibus-sat-ws/credential"
	"github.com/InsaneTreset/nibus-sat-ws/service"
)

func NewCFDIClient(cer, key []byte, pass string) (*service.Client, error) {
	cred, err := credential.NewFromBytes(cer, key, pass)
	if err != nil {
		return nil, err
	}
	if err := cred.Validate(); err != nil {
		return nil, err
	}
	return service.New(cred, service.CFDI), nil
}
