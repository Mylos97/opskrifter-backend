package types

type Recipe struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Minutes       int                `json:"minutes"`
	Description   string             `json:"description"`
	Likes         int                `json:"likes"`
	Comments      int                `json:"comments"`
	Image         string             `json:"image"`
	Ingredients   []IngredientAmount `json:"ingredients"`
	Categories    []RecipeCategory   `json:"categories"`
	RecipeCuisine RecipeCuisine      `json:"recipeCuisine"`
	User          User               `json:"user_id"`
}

type RecipeDB struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Minutes       int    `json:"minutes"`
	Description   string `json:"description"`
	Likes         int    `json:"likes"`
	Comments      int    `json:"comments"`
	Image         string `json:"image"`
	RecipeCuisine int    `json:"recipe_cuisine"`
	User          string `json:"user_id"`
}

func ToRecipeDB(r Recipe) RecipeDB {
	return RecipeDB{
		ID:            r.ID,
		Name:          r.Name,
		Minutes:       r.Minutes,
		Description:   r.Description,
		Likes:         r.Likes,
		Comments:      r.Comments,
		Image:         r.Image,
		RecipeCuisine: r.RecipeCuisine.ID,
		User:          r.User.ID,
	}
}

func FromRecipeDB(db RecipeDB, cuisine RecipeCuisine, user User) Recipe {
	return Recipe{
		ID:            db.ID,
		Name:          db.Name,
		Minutes:       db.Minutes,
		Description:   db.Description,
		Likes:         db.Likes,
		Comments:      db.Comments,
		Image:         db.Image,
		RecipeCuisine: cuisine,
		User:          user,
	}
}
