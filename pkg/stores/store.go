package stores

type Store struct {
	Id        string   `json:"id"`
	Name      string   `json:"name" validate:"required"`
	OwnerId   string   `json:"owner_id" validate:"required"`
	Type      string   `json:"type" validate:"required"`
	Location  Location `json:"location"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type StoreAvailability struct {
	Id           string `json:"id"`
	StoreId      string `json:"store_id" validate:"required"`
	Availability string `json:"availability" validate:"required"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type StoreRating struct {
	Id        string  `json:"id"`
	StoreId   string  `json:"store_id" validate:"required"`
	UserId    string  `json:"user_id" validate:"required"`
	Rating    float32 `json:"rating" validate:"required,min=1,max=5"`
	Message   string  `json:"message"`
	CreatedAt string  `json:"created_at"`
}

type StorePlan struct {
	Id        string  `json:"id"`
	StoreId   string  `json:"store_id" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Price     float32 `json:"price" validate:"required"`
	Currency  string  `json:"currency" validate:"required"`
	PlanType  string  `json:"plan_type" validate:"required"`
	Frequency string  `json:"frequency" validate:"required"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

type StoreAppointment struct {
	Id            string  `json:"id"`
	StoreId       string  `json:"store_id" validate:"required"`
	UserId        string  `json:"user_id" validate:"required"`
	StartAt       string  `json:"start_at" validate:"required"`
	EndAt         string  `json:"end_at" validate:"required"`
	Status        string  `json:"status" validate:"required"`
	HoldExpiresAt string  `json:"hold_expires_at"`
	Price         float32 `json:"price" validate:"required"`
	Currency      string  `json:"currency" validate:"required"`
	FeePlatform   float32 `json:"fee_platform" validate:"required"`
	PaymentId     string  `json:"payment_id"`
	Notes         string  `json:"notes"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

type Subscription struct {
	Id        string `json:"id"`
	StoreId   string `json:"store_id" validate:"required"`
	PlanId    string `json:"plan_id" validate:"required"`
	UserId    string `json:"user_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Location struct {
	Source      string      `json:"source" validate:"required"`
	Provider    string      `json:"provider" validate:"required"`
	Coordinates Coordinates `json:"coordinates" validate:"required"`
}

type Coordinates struct {
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Accuracy  float32 `json:"accuracy"` // esses dados podem vir diferente dependendo do provedor, verificar modos de ter fallback no json
}

type Availability struct {
	DayOfWeek string `json:"day_of_week" validate:"required"`
	OpenTime  string `json:"open_time" validate:"required"`
	CloseTime string `json:"close_time" validate:"required"`
}

type FullStore struct {
	Store        Store             `json:"store"`
	Availability StoreAvailability `json:"availability"`
	Ratings      []StoreRating     `json:"ratings"`
	Plans        []StorePlan       `json:"plans"`
}

type GetStoreResponse struct {
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Page   int     `json:"page"`
	Stores []Store `json:"stores"`
}

type GetStoreByIdResponse struct {
	Total int       `json:"total"`
	Store FullStore `json:"store"`
}

type GetStoreRatingsResponse struct {
	Total   int           `json:"total"`
	Limit   int           `json:"limit"`
	Page    int           `json:"page"`
	Ratings []StoreRating `json:"ratings"`
}

type GetStoreAppointmentsResponse struct {
	Total        int                `json:"total"`
	Limit        int                `json:"limit"`
	Page         int                `json:"page"`
	Appointments []StoreAppointment `json:"appointments"`
}
