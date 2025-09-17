package app

type Application interface {
	Run() error
	Close() error
}

var _ Application = (*App)(nil)
