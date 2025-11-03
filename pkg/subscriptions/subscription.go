package subscriptions

type Subscription struct {
	Id                 string `json:"id"`
	UserId             string `json:"user_id" validate:"required"`
	StoreId            string `json:"store_id" validate:"required"`
	StorePlanId        string `json:"store_plan_id" validate:"required"`
	Status             string `json:"status" validate:"required"`
	CurrentPeriodStart string `json:"current_period_start" validate:"required"`
	CurrentPeriodEnd   string `json:"current_period_end" validate:"required"`
	CancelAt           string `json:"cancel_at"`
	CanceledAt         string `json:"canceled_at"`
	TrialEnd           string `json:"trial_end"`
	Provider           string `json:"provider" validate:"required"`
	ProviderSubId      string `json:"provider_sub_id" validate:"required"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type GetSubscriptionsResponse struct {
	Total         int            `json:"total"`
	Limit         int            `json:"limit"`
	Subscriptions []Subscription `json:"subscriptions"`
}
