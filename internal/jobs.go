package internal

const (
	Completed string = "completed"
	Ongoing   string = "ongoing"
	Failed    string = "failed"
)

// Job will contian the information about the each job
type Job struct {
	//TODO: use 'count' field if required
	Count  int     `json:"count,omitempty"`
	JobID  int     `json:"job_id,omitempty"`
	Visits []Store `json:"visits,omitempty"`
	Status string  `json:"status"`
}

// Store is struct contains info about store
type Store struct {
	StoreID   string   `json:"store_id"`
	ImageURL  []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}
