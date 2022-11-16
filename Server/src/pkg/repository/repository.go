package repository

type Package interface {
	SavePackage(content []byte) error
}

type Repository struct {
	Package
}

func NewRepository() *Repository {
	return &Repository{}
}
