package api

import (
	"fmt"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"
)

func TestDeleteGeneric(t *testing.T) {
	fmt.Printf("%+v\n", testRecipe)
	id, err := CreateByType(testRecipe)

	if err != nil {
		t.Fatal(err)
	}

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	testRecipe.ID = id
	_, err = DeleteByType(testRecipe)
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

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	if err != nil {
		t.Fatalf("failed to insert recipe: %v", err)
	}

	testRecipe.ID = id
	obj, err := GetByType(testRecipe)

	if err != nil {
		t.Fatalf("failed to get recipe: %v", err)
	}

	if obj.GetID() != testRecipe.GetID() {
		t.Errorf("expected ID %s, got %s", testRecipe.GetID(), obj.GetID())
	}

	if obj.Name != testRecipe.Name {
		t.Errorf("expected Name %s, got %s", testRecipe.Name, obj.Name)
	}

	_, err = DeleteByType(obj)
	if err != nil {
		t.Fatalf("failed to clean up recipe: %v", err)
	}
}

func TestCreateGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	testRecipe.ID = id

	if err != nil {
		t.Fatalf("failed to insert recipe: %v", err)
	}
	err = testutils.AssertCountByType[types.Recipe](1, GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	_, err = DeleteByType(testRecipe)

	if err != nil {
		t.Fatalf("failed to clean up recipe: %v", err)
	}
}

func TestUpdateGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	if err != nil {
		t.Fatalf("failed to insert recipe for update test: %v", err)
	}

	testRecipe.Name = "Updated Recipe"
	testRecipe.Description = "Updated Description"
	testRecipe.Image = "after.jpg"
	testRecipe.Likes = 42
	testRecipe.ID = id

	_, err = UpdateByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to update recipe: %v", err)
	}

	updated, err := GetByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to fetch updated recipe: %v", updated)
	}

	if updated.Name != "Updated Recipe" || updated.Description != "Updated Description" ||
		updated.Image != "after.jpg" || updated.Likes != 42 {
		t.Errorf("update not applied correctly: %+v", updated)
	}

	_, err = DeleteByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to clean up recipe: %v", err)
	}
}
func TestGetMany(t *testing.T) {
	for i := range testRecipes {
		id, err := CreateByType(testRecipes[i])
		testRecipes[i].ID = id

		if err != nil {
			t.Fatalf("failed to insert recipe at index %d: %v\nRecipe ID: %s", i, err, id)
		}

		if id == "" {
			t.Fatalf("failed to create a id: %v", id)
		}

	}

	err := testutils.AssertCountByType[types.Recipe](len(testRecipes), GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	for i, recipe := range testRecipes {
		_, err := DeleteByType(recipe)
		if err != nil {
			t.Fatalf("failed to delete recipe at index %d: %v\nRecipe ID: %s", i, err, recipe.GetID())
		}
	}

	err = testutils.AssertCountByType[types.Recipe](0, GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}
}

func TestOneToMany(t *testing.T) {
	ids, err := CreateManyByType(testIngredients)

	if err != nil {
		t.Fatalf("error creating ingredients")
	}

	for i := range ids {
		testIngredients[i].ID = ids[i]
	}

	ids, err = CreateManyByType(testRecipes)
	if err != nil {
		t.Fatalf("error creating recipes")
	}

	for i := range ids {
		testRecipes[i].ID = ids[i]
	}

	for i := range testRecipes {
		recipeIngredients := types.ToOneToMany(
			testIngredients,
			testRecipes[i],
			types.IngredientToRecipeIngredient,
		)
		err = CreateOneToManyByType(testRecipes[i], testRecipes[i].ID, recipeIngredients)

		if err != nil {
			t.Fatalf("failed to insert relations at index %d: %v\n Recipe ID: %s", i, err, testRecipes[i].ID)
		}
	}
	tableName := types.RecipeIngredient{}.TableName()

	expectedLength := len(testIngredients) * len(testRecipes)
	err = testutils.AssertCountByTable(expectedLength, tableName, GetCountByTable)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	err = DeleteManyByType(testRecipes)

	if err != nil {
		t.Fatalf("error deleting recipes")
	}

	err = testutils.AssertCountByTable(0, tableName, GetCountByTable)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	err = DeleteManyByType(testIngredients)
	if err != nil {
		t.Fatalf("error deleting ingredients")
	}
}

func TestCreateByTypeWithRelations(t *testing.T) {
	ids, err := CreateManyByType(testIngredients)
	if err != nil {
		t.Fatalf("error creating ingredients")
	}

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

	if err != nil {
		t.Fatalf("error creating recipe with relations %v", err)
	}

	if id == "" {
		t.Fatalf("error generating a ID")
	}

	testRecipe.ID = id

	expectedLength := len(testIngredients)
	err = testutils.AssertCountByTable(expectedLength, tableName, GetCountByTable)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	_, err = DeleteByType(testRecipe)

	if err != nil {
		t.Fatalf("error deleting recipes %v", err)
	}

	err = testutils.AssertCountByType[types.Recipe](0, GetCountByType)

	if err != nil {
		t.Fatalf("%v", err)
	}

	err = testutils.AssertCountByTable(0, tableName, GetCountByTable)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}
}
