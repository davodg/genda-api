package subscriptions

type Service interface {
	CreateSubscription(Subscription) (*Subscription, error)
	GetSubscriptions(userId string, page int, limit int) (*GetSubscriptionsResponse, error)
	UpdateSubscription(id string, subscription Subscription) (*Subscription, error)
	DeleteSubscription(id string) error
}

type Repository interface {
	CreateSubscription(Subscription) (*Subscription, error)
	GetSubscriptions(userId string, page int, limit int) (*GetSubscriptionsResponse, error)
	UpdateSubscription(id string, subscription Subscription) (*Subscription, error)
	DeleteSubscription(id string) error
}

type service struct {
	subscriptionRepository Repository
}

func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) CreateSubscription(subscription Subscription) (*Subscription, error) {
	return s.subscriptionRepository.CreateSubscription(subscription)
}

func (s *service) GetSubscriptions(userId string, page int, limit int) (*GetSubscriptionsResponse, error) {
	return s.subscriptionRepository.GetSubscriptions(userId, page, limit)
}

func (s *service) UpdateSubscription(id string, subscription Subscription) (*Subscription, error) {
	return s.subscriptionRepository.UpdateSubscription(id, subscription)
}

func (s *service) DeleteSubscription(id string) error {
	return s.subscriptionRepository.DeleteSubscription(id)
}
