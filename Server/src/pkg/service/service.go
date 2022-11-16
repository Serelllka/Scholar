package service

type Package interface {
}

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
