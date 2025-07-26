package api

import (
	"opskrifter-backend/internal/types"
)

var CreateRecipe = HandlerByType(CreateByTypeWithRelations[types.Recipe])
var UpdateRecipe = HandlerByType(UpdateByType[types.Recipe])
var DeleteRecipe = HandlerByType(DeleteByType[types.Recipe])
var GetRecipe = GetHandlerByType(GetByType[types.Recipe])
