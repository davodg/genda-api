package users

type Service interface {
	CreateUser(User) (*User, error)
	GetUsers(int, int, string, string, string) (*GetUserResponse, error)
	GetUser(string) (*User, error)
	UpdateUser(string, User) (*User, error)
	DeleteUser(string) error
}

type Repository interface {
	CreateUser(User) (*User, error)
	GetUsers(int, int, string, string, string) (*GetUserResponse, error)
	GetUser(string) (*User, error)
	UpdateUser(string, User) (*User, error)
	DeleteUser(string) error
}

type service struct {
	userRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateUser(user User) (*User, error) {
	return s.userRepository.CreateUser(user)
}

func (s *service) GetUsers(page int, limit int, name string, email string, userId string) (*GetUserResponse, error) {
	return s.userRepository.GetUsers(page, limit, name, email, userId)
}

func (s *service) GetUser(id string) (*User, error) {
	res, err := s.userRepository.GetUsers(1, 1, "", "", id)
	return &res.Users[0], err
}

func (s *service) UpdateUser(id string, user User) (*User, error) {
	return s.userRepository.UpdateUser(id, user)
}

func (s *service) DeleteUser(id string) error {
	return s.userRepository.DeleteUser(id)
}
