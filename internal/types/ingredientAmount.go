package types

type IngredientAmount struct {
	ID         string     `json:"id"`
	Ingredient Ingredient `json:"name"`
	Amount     string     `json:"amount"`
	Recipe     Recipe     `json:"recipe"`
}
