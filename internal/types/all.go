package types

// Recipe
type Recipe struct {
	ID                string `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Minutes           int    `json:"minutes" db:"minutes"`
	Description       string `json:"description" db:"description"`
	Likes             int    `json:"likes" db:"likes"`
	Comments          int    `json:"comments" db:"comments"`
	Views             int    `json:"views" db:"views"`
	Image             string `json:"image" db:"image"`
	RecipeCuisine     string `json:"recipe_cuisine" db:"recipe_cuisine"`
	UserID            string `json:"user_id" db:"user_id"`
	CreatedAt         string `json:"created_at" db:"created_at"`
	RecipeIngredients []RecipeIngredient
}

func (Recipe) TableName() string { return "recipes" }
func (r Recipe) GetID() string   { return r.ID }
func (r Recipe) GetManyToMany() [][]ManyToMany {
	parts := [][]ManyToMany{}
	parts = append(parts, ToInterfaceSlice(r.RecipeIngredients))
	return parts
}

// User
type User struct {
	ID        string `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	Email     string `json:"email" db:"email"`
	Status    string `json:"status" db:"status"`
	CreatedAt string `json:"created_at" db:"created_at"`
}

func (User) TableName() string { return "users" }
func (u User) GetID() string   { return u.ID }

// UserLikedRecipe
type UserLikedRecipe struct {
	UserID   string `json:"user_id" db:"user_id" parent:"true"`
	RecipeID string `json:"recipe_id" db:"recipe_id" child:"true"`
}

func (UserLikedRecipe) TableName() string      { return "user_liked_recipes" }
func (ulr UserLikedRecipe) GetChildID() string { return ulr.RecipeID }

// Ingredient
type Ingredient struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

func (Ingredient) TableName() string { return "ingredients" }
func (i Ingredient) GetID() string   { return i.ID }

// RecipeIngredient
type RecipeIngredient struct {
	RecipeId     string `json:"recipe_id" db:"recipe_id" parent:"true"`
	IngredientId string `json:"ingredient_id" db:"ingredient_id" child:"true"`
	Amount       string `json:"amount" db:"amount"`
}

func (RecipeIngredient) TableName() string     { return "ingredients_for_recipe" }
func (ri RecipeIngredient) GetChildID() string { return ri.IngredientId }
