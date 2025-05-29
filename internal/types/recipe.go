package types

type Recipe struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Minutes      int                    `json:"minutes"`
	Rating       float64                `json:"rating"`
	Description  string                 `json:"description"`
	Likes        int                    `json:"likes"`
	Comments     int                    `json:"comments"`
	Image        string           			`json:"image"`
	Ingredients  []IngredientAmount     `json:"ingredients"`
	Categories   []RecipeCategory 			`json:"categories"`
	RecipeCuisine RecipeCuisine   			`json:"recipeCuisine"`
	User				 User										`json:"user"`
}
