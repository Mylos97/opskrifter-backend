package api

import (
	"fmt"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"
)

var recipeGenerator = testutils.NewTestDataGenerator[types.RecipeV2]()
var testRecipe = recipeGenerator.Generate()
var amount = 1000
var testRecipes = recipeGenerator.GenerateMany(amount)

func TestDeleteGeneric(t *testing.T) {
	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), testRecipe.ID)

	err := CreateByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to insert recipe: %v", err)
	}

	err = DeleteByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to delete recipe: %v", err)
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
	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), testRecipe.ID)

	err := CreateByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to insert recipe: %v", err)
	}

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

	err = DeleteByType(obj)
	if err != nil {
		t.Fatalf("failed to clean up recipe: %v", err)
	}
}

func TestCreateGeneric(t *testing.T) {
	fmt.Printf("Inserting into table %s with ID: %s\n", testRecipe.TableName(), testRecipe.ID)

	err := CreateByType(testRecipe)
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

	_ = DeleteByType(testRecipe)
}

func TestUpdateGeneric(t *testing.T) {
	err := CreateByType(testRecipe)
	if err != nil {
		t.Fatalf("failed to insert recipe for update test: %v", err)
	}

	testRecipe.Name = "Updated Recipe"
	testRecipe.Description = "Updated Description"
	testRecipe.Image = "after.jpg"
	testRecipe.Likes = 42

	err = UpdateByType(testRecipe)
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

	_ = DeleteByType(testRecipe)
}

func TestGetMany(t *testing.T) {
	for i, recipe := range testRecipes {
		err := CreateByType(recipe)
		if err != nil {
			t.Fatalf("failed to insert recipe at index %d: %v\nRecipe ID: %s", i, err, recipe.GetID())
		}
	}

	count, err := GetCountByType(testRecipe)

	if err != nil {
		t.Fatalf("error getting the count")
	}

	if count != amount {
		t.Fatalf("Expecting %d got %d", amount, count)
	}

	for i, recipe := range testRecipes {
		err := DeleteByType(recipe)
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
