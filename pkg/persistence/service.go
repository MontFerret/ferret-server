package persistence

import "github.com/MontFerret/ferret-server/pkg/scripts"

type Service struct {
}

func (service *Service) Collect(script scripts.Script, data []byte) error {
	return nil
}
