package documentmodels

type Collection struct {
	Name      string     `json:"name"`
	Documents []Document `json:"documents"`
}
