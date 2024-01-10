package listing

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Order          int    `json:"order"`
	BaseExperience int    `json:"base_experience"`
	IsDefault      bool   `json:"is_default"`
}
