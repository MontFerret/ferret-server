package execution

import "runtime"

type Settings struct {
	PoolSize uint64
}

func NewDefaultSettings() Settings {
	return Settings{
		PoolSize: uint64(runtime.NumCPU() * 10),
	}
}
