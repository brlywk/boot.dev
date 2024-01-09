package routes

// ----- Types -----------------------------------

type RequestBody struct {
	Body string `json:"body"`
}

type CleanResponse struct {
	CleanedBody string `json:"cleaned_body"`
}

type MetricsConfig struct {
	FileserverHits int
}
