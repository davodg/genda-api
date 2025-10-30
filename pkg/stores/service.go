package stores

type Service interface {
	CreateStore(Store) (*Store, error)
	GetStores(int, int, string, string, string) (*GetStoreResponse, error)
	GetStore(string) (*Store, error)
	UpdateStore(string, Store) (*Store, error)
	DeleteStore(string) error
}

type Repository interface {
	CreateStore(Store) (*Store, error)
	GetStores(int, int, string, string, string) (*GetStoreResponse, error)
	GetStore(string) (*Store, error)
	UpdateStore(string, Store) (*Store, error)
	DeleteStore(string) error
}

type service struct {
	storeRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateStore(store Store) (*Store, error) {
	return s.storeRepository.CreateStore(store)
}

func (s *service) GetStores(page int, limit int, name string, storeType string, storeId string) (*GetStoreResponse, error) {
	return s.storeRepository.GetStores(page, limit, name, storeType, storeId)
}

func (s *service) GetStore(id string) (*Store, error) {
	res, err := s.storeRepository.GetStores(1, 1, "", "", id)
	return &res.Stores[0], err
}

func (s *service) UpdateStore(id string, store Store) (*Store, error) {
	return s.storeRepository.UpdateStore(id, store)
}

func (s *service) DeleteStore(id string) error {
	return s.storeRepository.DeleteStore(id)
}
