package api

import (
	"log"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/myDB"
	"os"
	"testing"
)

var adminUser = types.User{
	Name:      "Test Admin",
	Email:     "admin@example.com",
	Status:    "admin",
	CreatedAt: "now",
}

var (
	recipeGenerator     *testutils.TestDataGenerator[types.Recipe]
	ingredientGenerator *testutils.TestDataGenerator[types.Ingredient]
	testRecipe          types.Recipe
	handlerRecipe       types.Recipe
	testRecipes         []types.Recipe
	testIngredients     []types.Ingredient
	amount              int
)

func TestMain(m *testing.M) {
	err := myDB.Init(true)
	if err != nil {
		log.Fatalf("error init DB %s", err)
	}

	id, err := CreateByType(adminUser)
	if err != nil {
		log.Fatal("error creating admin user")
	}
	adminUser.ID = id

	setupTestData()
	code := m.Run()
	myDB.DB.Close()
	os.Exit(code)
}

func setupTestData() {
	amount = 10
	recipeGenerator = testutils.NewTestDataGenerator[types.Recipe](adminUser.ID)
	ingredientGenerator = testutils.NewTestDataGenerator[types.Ingredient](adminUser.ID)

	testRecipe = recipeGenerator.Generate()
	testRecipes = recipeGenerator.GenerateMany(amount)
	handlerRecipe = recipeGenerator.Generate()

	testIngredients = ingredientGenerator.GenerateMany(amount)
}
