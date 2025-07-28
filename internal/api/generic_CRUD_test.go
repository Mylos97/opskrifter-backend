package api

import (
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)
	require.NotEmpty(t, id, "failed to create an ID")
	require.NoError(t, err, "failed to insert recipe for delete test")

	testRecipe.ID = id
	_, err = DeleteByType[types.Recipe](id)
	if err != nil {
		t.Fatal(err)
	}

	err = testutils.AssertCountByType[types.Recipe](0, GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}
}

func TestGetGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)

	require.NotEmpty(t, id, "failed to create an ID")
	require.NoError(t, err, "failed to insert recipe for get test")

	testRecipe.ID = id
	obj, err := GetByType(testRecipe)

	require.NoError(t, err, "failed to get recipe")
	assert.Equal(t, testRecipe.GetID(), obj.GetID(), "unexpected recipe ID")
	assert.Equal(t, testRecipe.Name, obj.Name, "unexpected recipe Name")

	_, err = DeleteByType[types.Recipe](id)
	if err != nil {
		t.Fatalf("failed to clean up recipe: %v", err)
	}
}

func TestCreateGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)
	require.NoError(t, err, "failed to insert recipe")
	require.NotEmpty(t, id, "failed to create an ID")
	err = testutils.AssertCountByType[types.Recipe](1, GetCountByType)
	require.NoError(t, err, "failed to get the count")

	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "failed to clean up recipe")
}

func TestUpdateGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)

	require.NotEmpty(t, id, "failed to create an ID")
	require.NoError(t, err, "failed to insert recipe for update test")

	require.NoError(t, err, "failed to insert recipe for update test")

	testRecipe.Name = "Updated Recipe"
	testRecipe.Description = "Updated Description"
	testRecipe.Image = "after.jpg"
	testRecipe.Likes = 42
	testRecipe.ID = id

	_, err = UpdateByType(testRecipe)
	require.NoError(t, err, "failed to update recipe")

	updated, err := GetByType(testRecipe)
	require.NoError(t, err, "failed to fetch updated recipe")

	assert.Equal(t, "Updated Recipe", updated.Name, "name was not updated correctly")
	assert.Equal(t, "Updated Description", updated.Description, "description was not updated correctly")
	assert.Equal(t, "after.jpg", updated.Image, "image was not updated correctly")
	assert.Equal(t, 42, updated.Likes, "likes was not updated correctly")

	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "failed to clean up recipe")
}

func TestGetMany(t *testing.T) {
	for i := range testRecipes {
		id, err := CreateByType(testRecipes[i])
		testRecipes[i].ID = id
		require.NoErrorf(t, err, "failed to insert recipe at index %d", i)
		require.NotEmptyf(t, id, "failed to create an ID at index %d", i)
	}

	require.NoError(t, testutils.AssertCountByType[types.Recipe](len(testRecipes), GetCountByType), "failed to get the count")

	for i, recipe := range testRecipes {
		_, err := DeleteByType[types.Recipe](testRecipes[i].ID)
		require.NoErrorf(t, err, "failed to delete recipe at index %d (ID: %s)", i, recipe.GetID())
	}

	require.NoError(t, testutils.AssertCountByType[types.Recipe](0, GetCountByType), "failed to get the count after deletions")

}

func TestOneToMany(t *testing.T) {
	ingredientIDs, err := CreateManyByType(testIngredients)
	require.NoError(t, err, "error creating ingredients")
	for i := range ingredientIDs {
		testIngredients[i].ID = ingredientIDs[i]
	}

	recipeIDs, err := CreateManyByType(testRecipes)
	require.NoError(t, err, "error creating recipes")
	for i := range recipeIDs {
		testRecipes[i].ID = recipeIDs[i]
	}

	for i := range testRecipes {
		recipeIngredients := types.ToOneToMany(
			testIngredients,
			testRecipes[i],
			types.IngredientToRecipeIngredient,
		)
		err = CreateOneToManyByType(testRecipes[i], testRecipes[i].ID, recipeIngredients)
		require.NoErrorf(t, err, "failed to insert relations at index %d (Recipe ID: %s)", i, testRecipes[i].ID)
	}

	tableName := types.RecipeIngredient{}.TableName()
	expectedLength := len(testIngredients) * len(testRecipes)

	require.NoError(t, testutils.AssertCountByTable(expectedLength, tableName, GetCountByTable), "failed to get the count")
	require.NoError(t, DeleteManyByType[types.Recipe](recipeIDs), "error deleting recipes")
	require.NoError(t, testutils.AssertCountByTable(0, tableName, GetCountByTable), "failed to get the count after deleting recipes")
	require.NoError(t, DeleteManyByType[types.Ingredient](ingredientIDs), "error deleting ingredients")
}

func TestCreateByTypeWithRelations(t *testing.T) {
	ids, err := CreateManyByType(testIngredients)
	require.NoError(t, err, "error creating ingredients")

	for i := range ids {
		testIngredients[i].ID = ids[i]
	}

	recipeIngredients := types.ToOneToMany(
		testIngredients,
		testRecipe,
		types.IngredientToRecipeIngredient,
	)

	testRecipe.RecipeIngredients = recipeIngredients
	id, err := CreateByTypeWithRelations(testRecipe)
	tableName := types.RecipeIngredient{}.TableName()

	require.NoError(t, err, "error creating recipe with relations")
	require.NotEmpty(t, id, "error generating a ID")

	testRecipe.ID = id

	expectedLength := len(testIngredients)
	require.NoError(t, testutils.AssertCountByTable(expectedLength, tableName, GetCountByTable), "failed to get the count")

	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "error deleting recipes")

	require.NoError(t, testutils.AssertCountByType[types.Recipe](0, GetCountByType))
	require.NoError(t, testutils.AssertCountByTable(0, tableName, GetCountByTable), "failed to get the count after deletions")
}
