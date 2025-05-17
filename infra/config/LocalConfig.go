package config

type LocalConfig struct {
}

func (c LocalConfig) UseFzf() bool {
	return false
}
