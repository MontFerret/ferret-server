package persistence

import "github.com/MontFerret/ferret-server/pkg/scripts"

type Service struct {
}

func (service *Service) Collect(_ scripts.Script, _ []byte) error {
	return nil
}
