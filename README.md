## Endpoints

### 📘 Recipes

| Method | Endpoint           | Description              |
|--------|--------------------|--------------------------|
| POST   | `/recipes/`        | Create a new recipe      |
| GET    | `/recipes/{id}`    | Get a recipe by ID       |
| PUT    | `/recipes/{id}`    | Update a recipe by ID    |
| DELETE | `/recipes/{id}`    | Delete a recipe by ID    |
| GET    | `/recipes/`        | Get a list of recipes    |

<pre lang="md">
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "Spaghetti Carbonara",
  "minutes": 30,
  "description": "A classic Italian pasta dish made with eggs, ....",
  "likes": 25,
  "comments": 4,
  "image": "https://example.com/images/spaghetti.jpg",
  "ingredients": [
    {
      "ingredient": {
        "id": "ing-1",
        "name": "Spaghetti"
      },
      "amount": "200g"
    },
  ],
  "categories": [
    {
      "id": "cat-1",
      "name": "Pasta"
    }
  ],
  "recipe_cuisine": {
    "id": "cuisine-1",
    "name": "Italian"
  },
  "user_id": {
    "id": "user-123",
    "name": "Jane Doe",
    "email": "jane@example.com"
  }
}
</pre>

### 📚 Cookbooks

| Method | Endpoint             | Description                |
|--------|----------------------|----------------------------|
| POST   | `/cookbooks/`        | Create a new cookbook      |
| GET    | `/cookbooks/{id}`    | Get a cookbook by ID       |
| PUT    | `/cookbooks/{id}`    | Update a cookbook by ID    |
| DELETE | `/cookbooks/{id}`    | Delete a cookbook by ID    |

<pre lang="md">
{
  "id": "cb-001",
  "name": "Italian Favorites",
  "description": "A collection of classic Italian recipes passed down through generations.",
  "likes": 42,
  "recipes": [
    {
      ...
    },
    {
      ...
    }
  ],
  "creator": "user-123"
}
</pre>
### 💬 Comments

| Method | Endpoint                  | Description                         |
|--------|---------------------------|-------------------------------------|
| POST   | `/comments/`              | Add a new comment to a recipe       |
| GET    | `/comments/{recipe_id}`   | Get all comments for a recipe       |
| PUT    | `/comments/{id}`          | Update a comment by ID              |
| DELETE | `/comments/{id}`          | Delete a comment by ID              |

<pre lang="md">
{
  "id": "comment-001",
  "recipe": {
    ...
  },
  "user": {
    ...
  },
  "comment": "This recipe was amazing! I added a bit more garlic and it turned out great."
}
</pre>

### 👤 Users

| Method | Endpoint        | Description              |
|--------|-----------------|--------------------------|
| POST   | `/users/`       | Create a new user        |
| GET    | `/users/{id}`   | Get a user by ID         |
| PUT    | `/users/{id}`   | Update a user by ID      |
| DELETE | `/users/{id}`   | Delete a user by ID      |

<pre lang="md">
{
  "id": "user-123",
  "name": "Jane Doe",
  "email": "jane.doe@example.com",
  "createdAt": "2025-06-01T14:30:00Z"
}
</pre>

### ❤️ Recipe Likes

| Method | Endpoint                    | Description                    |
|--------|-----------------------------|--------------------------------|
| PUT    | `/like_recipe/like`         | Like a recipe                  |
| PUT    | `/like_recipe/unlike`       | Unlike a recipe                |
| GET    | `/like_recipe/{user_id}`  | Get the recipes a user liked   |

<pre lang="md">
{
  "user_id": "user-123",
  "recipe_id": "recipe-456"
}
</pre>