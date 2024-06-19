package controller_interfaces

type FindAllResponse[T any] struct {
	Total     uint32
	Page      uint16
	Registers []*T
}

type IGenericController[T any] interface {
	Create(m *T) (result string, err error)
	FindBy(id string) (result *T, err error)
	FindAll(page uint16) (result *FindAllResponse[T], err error)
	UpdateBy(id string, m *T) (result bool, err error)
	DeleteBy(id string) (result bool, err error)
}
