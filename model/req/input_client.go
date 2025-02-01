package req

type NewClient struct {
	Email   string `json:"email"`
	Balance int    `json:"balance"`
	Plan    string `json:"plan"`
}