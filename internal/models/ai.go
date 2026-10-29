package models

type FailureRequest struct {
	StatusCode   int    `json:"status_code"`
	ResponseBody string `json:"response_body"`
	Attempt      int    `json:"attempt"`
}

type FailureResponse struct {
	LikelyCause       string   `json:"likely_cause"`
	Evidence          []string `json:"evidence"`
	RecommendedAction string   `json:"recommended_action"`
	Confidence        string   `json:"confidence"`
}

type InvestigationRequest struct {
	Delivery Delivery        `json:"delivery"`
	Attempts []DeliveryAttempt `json:"attempts"`
}

