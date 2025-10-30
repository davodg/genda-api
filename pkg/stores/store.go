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

type GetStoreResponse struct {
	Total  int     `json:"total"`
	Limit  int     `json:"limit"`
	Page   int     `json:"page"`
	Stores []Store `json:"stores"`
}
