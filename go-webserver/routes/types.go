package routes

// ----- Types -----------------------------------

type ChirpRequestBody struct {
	Body string `json:"body"`
}

type UserRequestBody struct {
	Email string `json:"email"`
}

type MetricsConfig struct {
	FileserverHits int
}
