package types

type RecipeV2 struct {
	ID            string `json:"id" db:"id"`
	Name          string `json:"name" db:"name"`
	Minutes       int    `json:"minutes" db:"minutes"`
	Description   string `json:"description" db:"description"`
	Likes         int    `json:"likes" db:"likes"`
	Comments      int    `json:"comments" db:"comments"`
	Image         string `json:"image" db:"image"`
	RecipeCuisine string `json:"recipe_cuisine" db:"recipe_cuisine"`
	UserID        string `json:"user_id" db:"user_id"`
}

func (RecipeV2) TableName() string { return "recipes" }

func (r RecipeV2) GetID() string {
	return r.ID
}
