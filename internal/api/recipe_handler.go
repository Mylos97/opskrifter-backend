package api

import (
	"opskrifter-backend/internal/types"
)

var CreateRecipe = HandlerByType(CreateByType[types.RecipeV2])
var UpdateRecipe = HandlerByType(UpdateByType[types.RecipeV2])
var DeleteRecipe = HandlerByType(DeleteByType[types.RecipeV2])
var GetRecipe = GetHandlerByType(GetByType[types.RecipeV2])
