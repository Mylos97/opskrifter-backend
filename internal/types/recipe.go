package types

type Recipe struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Minutes       int    `json:"minutes" db:"minutes"`
	Description   string `json:"description" db:"description"`
	Likes         int    `json:"likes" db:"likes"`
	Comments      int    `json:"comments" db:"comments"`
	Image         string `json:"image" db:"image"`
	RecipeCuisine string `json:"recipe_cuisine" db:"recipe_cuisine"`
	UserID        string `json:"user_id" db:"user_id"`
	Ingredients   []RecipeIngredient
}

func (Recipe) TableName() string { return "recipes" }

func (r Recipe) GetID() string {
	return r.ID
}

func (r Recipe) GetOneToMany() [][]OneToMany {
	parts := [][]OneToMany{}
	return parts
}

type Ingredient struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type RecipeIngredient struct {
	ID           string `json:"id" db:"id"`
	RecipeId     string `json:"recipe_id" db:"recipe_id" parent:""`
	IngredientId string `json:"ingredient_id" db:"ingredient_id" child:""`
}
