package stores

type Service interface {
	CreateStore(Store) (*Store, error)
	GetStores(int, int, string, string, string) (*GetStoreResponse, error)
	GetStore(string) (*GetStoreByIdResponse, error)
	UpdateStore(string, Store) (*Store, error)
	DeleteStore(string) error
	CreateStorePlan(StorePlan) (*StorePlan, error)
	CreateStoreAvailability(StoreAvailability) (*StoreAvailability, error)
	CreateStoreRating(StoreRating) (*StoreRating, error)
	GetStoreRatings(string, int, int) (*GetStoreRatingsResponse, error)
	GetStoreAppointments(string, int, int) (*GetStoreAppointmentsResponse, error)
	UpdateStoreRating(string, StoreRating) (*StoreRating, error)
	DeleteStoreRating(string) error
	UpdateStorePlan(string, StorePlan) (*StorePlan, error)
	DeleteStorePlan(string) error
	UpdateStoreAvailability(string, StoreAvailability) (*StoreAvailability, error)
	DeleteStoreAvailability(string) error
	GetStorePlans(string) (*[]StorePlan, error)
	GetStoreAvailability(string) (*StoreAvailability, error)
	CreateStoreAppointment(StoreAppointment) (*StoreAppointment, error)
	UpdateStoreAppointment(string, StoreAppointment) (*StoreAppointment, error)
	DeleteStoreAppointment(string) error
}

type Repository interface {
	CreateStore(Store) (*Store, error)
	GetStores(int, int, string, string, string) (*GetStoreResponse, error)
	GetStore(string) (*GetStoreByIdResponse, error)
	UpdateStore(string, Store) (*Store, error)
	DeleteStore(string) error
	CreateStorePlan(StorePlan) (*StorePlan, error)
	CreateStoreAvailability(StoreAvailability) (*StoreAvailability, error)
	CreateStoreRating(StoreRating) (*StoreRating, error)
	GetStoreRatings(string, int, int) (*GetStoreRatingsResponse, error)
	GetStoreAppointments(string, int, int) (*GetStoreAppointmentsResponse, error)
	UpdateStoreRating(string, StoreRating) (*StoreRating, error)
	DeleteStoreRating(string) error
	UpdateStorePlan(string, StorePlan) (*StorePlan, error)
	DeleteStorePlan(string) error
	UpdateStoreAvailability(string, StoreAvailability) (*StoreAvailability, error)
	DeleteStoreAvailability(string) error
	GetStorePlans(string) (*[]StorePlan, error)
	GetStoreAvailability(string) (*StoreAvailability, error)
	CreateStoreAppointment(StoreAppointment) (*StoreAppointment, error)
	UpdateStoreAppointment(string, StoreAppointment) (*StoreAppointment, error)
	DeleteStoreAppointment(string) error
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

func (s *service) GetStore(id string) (*GetStoreByIdResponse, error) {
	res, err := s.storeRepository.GetStore(id)
	return res, err
}

func (s *service) UpdateStore(id string, store Store) (*Store, error) {
	return s.storeRepository.UpdateStore(id, store)
}

func (s *service) DeleteStore(id string) error {
	return s.storeRepository.DeleteStore(id)
}

func (s *service) CreateStorePlan(plan StorePlan) (*StorePlan, error) {
	return s.storeRepository.CreateStorePlan(plan)
}

func (s *service) CreateStoreAvailability(availability StoreAvailability) (*StoreAvailability, error) {
	return s.storeRepository.CreateStoreAvailability(availability)
}

func (s *service) CreateStoreRating(rating StoreRating) (*StoreRating, error) {
	return s.storeRepository.CreateStoreRating(rating)
}

func (s *service) GetStoreRatings(storeId string, page int, limit int) (*GetStoreRatingsResponse, error) {
	return s.storeRepository.GetStoreRatings(storeId, page, limit)
}

func (s *service) GetStoreAppointments(storeId string, page int, limit int) (*GetStoreAppointmentsResponse, error) {
	return s.storeRepository.GetStoreAppointments(storeId, page, limit)
}

func (s *service) UpdateStoreRating(id string, rating StoreRating) (*StoreRating, error) {
	return s.storeRepository.UpdateStoreRating(id, rating)
}

func (s *service) DeleteStoreRating(id string) error {
	return s.storeRepository.DeleteStoreRating(id)
}

func (s *service) UpdateStorePlan(id string, plan StorePlan) (*StorePlan, error) {
	return s.storeRepository.UpdateStorePlan(id, plan)
}

func (s *service) DeleteStorePlan(id string) error {
	return s.storeRepository.DeleteStorePlan(id)
}

func (s *service) UpdateStoreAvailability(id string, availability StoreAvailability) (*StoreAvailability, error) {
	return s.storeRepository.UpdateStoreAvailability(id, availability)
}

func (s *service) DeleteStoreAvailability(id string) error {
	return s.storeRepository.DeleteStoreAvailability(id)
}

func (s *service) GetStorePlans(storeId string) (*[]StorePlan, error) {
	return s.storeRepository.GetStorePlans(storeId)
}

func (s *service) GetStoreAvailability(storeId string) (*StoreAvailability, error) {
	return s.storeRepository.GetStoreAvailability(storeId)
}

func (s *service) CreateStoreAppointment(appointment StoreAppointment) (*StoreAppointment, error) {
	return s.storeRepository.CreateStoreAppointment(appointment)
}

func (s *service) UpdateStoreAppointment(id string, appointment StoreAppointment) (*StoreAppointment, error) {
	return s.storeRepository.UpdateStoreAppointment(id, appointment)
}

func (s *service) DeleteStoreAppointment(id string) error {
	return s.storeRepository.DeleteStoreAppointment(id)
}
