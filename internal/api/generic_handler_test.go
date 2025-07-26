package api

import (
	"fmt"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"
)

var amount = 10
var recipeGenerator = testutils.NewTestDataGenerator[types.Recipe]()
var ingredientGenerator = testutils.NewTestDataGenerator[types.Ingredient]()
var testRecipe = recipeGenerator.Generate()
var testRecipes = recipeGenerator.GenerateMany(amount)
var testIngredients = ingredientGenerator.GenerateMany(amount)

func TestDeleteGeneric(t *testing.T) {

	id, err := CreateByType(testRecipe)

	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), id)

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	if err != nil {
		t.Fatal(err)
	}
	testRecipe.ID = id
	_, err = DeleteByType(testRecipe)
	if err != nil {
		t.Fatal(err)
	}

	var count = 0
	count, err = GetCountByType(testRecipe)

	if err != nil {
		t.Fatalf("failed to verify deletion: %v", err)
	}
	if count != 0 {
		t.Errorf("expected 0 rows after delete, found %d", count)
	}
}

func TestGetGeneric(t *testing.T) {
	id, err := CreateByType(testRecipe)

	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), testRecipe.ID)

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
	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), testRecipe.ID)

	id, err := CreateByType(testRecipe)

	if id == "" {
		t.Fatalf("failed to create a id: %v", id)
	}

	testRecipe.ID = id

	if err != nil {
		t.Fatalf("failed to insert recipe: %v", err)
	}
	var count = 0
	count, err = GetCountByType(testRecipe)

	if err != nil {
		t.Fatalf("failed to verify insert: %v", err)
	}

	if count != 1 {
		t.Errorf("expected 1 row after insert")
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

	count, err := GetCountByType(testRecipe)

	if err != nil {
		t.Fatalf("error getting the count")
	}

	if count != len(testRecipes) {
		t.Fatalf("Expecting %d got %d", amount, count)
	}

	for i, recipe := range testRecipes {
		_, err := DeleteByType(recipe)
		if err != nil {
			t.Fatalf("failed to delete recipe at index %d: %v\nRecipe ID: %s", i, err, recipe.GetID())
		}
	}

	count, err = GetCountByType(testRecipe)

	if err != nil {
		t.Fatalf("error getting the count")
	}

	if count != 0 {
		t.Fatalf("Expecting 0 got %d", count)
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
		err = CreateOneToManyByType(testRecipes[i], recipeIngredients)

		if err != nil {
			t.Fatalf("failed to insert relations at index %d: %v\n Recipe ID: %s", i, err, testRecipes[i].ID)
		}
	}
	tableName := types.RecipeIngredient{}.TableName()
	expectedLength := len(testIngredients) * len(testRecipes)
	count, err := GetCountByTable(tableName)

	if err != nil {
		t.Fatalf("error getting the count")
	}

	if count != expectedLength {
		t.Fatalf("Expecting %d got %d", expectedLength, count)
	}
	println(len(testRecipes))
	err = DeleteManyByType(testRecipes)

	if err != nil {
		t.Fatalf("error deleting recipes")
	}

	count, _ = GetCountByType(testRecipes[0])

	if count != 0 {
		t.Fatalf("expected 0 got %d", count)
	}

	count, err = GetCountByTable(tableName)

	if err != nil {
		t.Fatalf("error getting the count")
	}

	if count != 0 {
		t.Fatalf("expected 0 got %d", count)
	}
}
