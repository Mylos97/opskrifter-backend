package types

type Recipe struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Minutes       int                    `json:"minutes"`
	Description   string                 `json:"description"`
	Likes         int                    `json:"likes"`
	Comments      int                    `json:"comments"`
	Image         string                 `json:"image"`
	Ingredients   []IngredientsForRecipe `json:"ingredients"`
	Categories    []RecipeCategory       `json:"categories"`
	RecipeCuisine string                 `json:"recipe_cuisine"`
	User          User                   `json:"user_id"`
}
