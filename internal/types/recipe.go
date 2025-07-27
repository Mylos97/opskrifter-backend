package types

type Recipe struct {
	ID                string `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Minutes           int    `json:"minutes" db:"minutes"`
	Description       string `json:"description" db:"description"`
	Likes             int    `json:"likes" db:"likes"`
	Comments          int    `json:"comments" db:"comments"`
	Image             string `json:"image" db:"image"`
	RecipeCuisine     string `json:"recipe_cuisine" db:"recipe_cuisine"`
	UserID            string `json:"user_id" db:"user_id"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	RecipeIngredients []RecipeIngredient
}

func (Recipe) TableName() string { return "recipes" }

func (r Recipe) GetID() string { return r.ID }

func (r Recipe) GetOneToMany() [][]OneToMany {
	parts := [][]OneToMany{}
	parts = append(parts, ToInterfaceSlice(r.RecipeIngredients))
	return parts
}

type Ingredient struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (Ingredient) TableName() string { return "ingredients" }

func (i Ingredient) GetID() string { return i.ID }

type RecipeIngredient struct {
	ID           string `json:"id" db:"id"`
	RecipeId     string `json:"recipe_id" db:"recipe_id" parent:"true"`
	IngredientId string `json:"ingredient_id" db:"ingredient_id" child:"true"`
}

func (RecipeIngredient) TableName() string { return "ingredients_for_recipe" }

func (ri RecipeIngredient) GetChildID() string {
	return ri.IngredientId
}
