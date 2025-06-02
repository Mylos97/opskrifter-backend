package types

type UserLikedRecipe struct {
	UserId   string `json:"user_id"`
	RecipeId string `json:"recipe_id"`
}
