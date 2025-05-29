package types

type CookBook struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Likes       int      `json:"likes"`
	Recipes     []Recipe `json:"recipes"`
}