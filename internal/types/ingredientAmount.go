package types

type IngredientAmount struct {
	ID         string     `json:"id"`
	Ingredient Ingredient `json:"name"`
	Quantity   int64      `json:"quantity"`
	Recipe     Recipe     `json:"recipe"`
}
