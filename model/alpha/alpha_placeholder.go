package alpha

type (
	AlphaRequest struct {
		PlaceholderID string `json:"placeholder_id"`
		Amount        int    `json:"amount"`
	}

	AlphaResponse struct {
		ID            string `json:"id"`
		PlaceholderID string `json:"placeholder_id"`
		Amount        int    `json:"amount"`
		Status        string `json:"status"`
	}
)
