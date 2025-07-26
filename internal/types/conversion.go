package types

func ToInterfaceSlice[T OneToMany](slice []T) []OneToMany {
	result := make([]OneToMany, len(slice))
	for i, v := range slice {
		result[i] = v
	}
	return result
}

func ToOneToMany[T Identifiable, R OneToMany](elements []T, parent Identifiable, factory func(T, Identifiable) R) []R {
	result := make([]R, 0, len(elements))
	for _, e := range elements {
		result = append(result, factory(e, parent))
	}
	return result
}

func IngredientToRecipeIngredient(ing Ingredient, rec Identifiable) RecipeIngredient {
	return RecipeIngredient{
		RecipeId:     rec.GetID(),
		IngredientId: ing.GetID(),
	}
}
