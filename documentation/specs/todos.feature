Feature: Todo Management
  As a user
  I want to be able to manage my todos
  So that I can keep track of my tasks

  Scenario: Create a new todo
    Given the backend server is running
    When I send a POST request to "/api/todos" with the following JSON:
      """
      {
        "title": "Buy groceries",
        "content": "Milk, eggs, bread",
        "stage": "To Do"
      }
      """
    Then the response status code should be 201
    And the response should contain a todo with title "Buy groceries", content "Milk, eggs, bread", and stage "To Do"

  Scenario: Retrieve all todos
    Given the database contains a todo with title "Buy groceries"
    When I send a GET request to "/api/todos"
    Then the response status code should be 200
    And the response should contain a list of todos
    And the list should include a todo with title "Buy groceries"

  Scenario: Update a todo
    Given the database contains a todo with title "Buy groceries" and stage "To Do"
    When I send a PUT request to update the todo to:
      """
      {
        "title": "Buy groceries",
        "stage": "In Progress"
      }
      """
    Then the response status code should be 200
    And the updated todo should have the stage "In Progress"

  Scenario: Toggle a todo completion status
    Given the database contains an incomplete todo with title "Buy groceries"
    When I send a PATCH request to toggle the todo completion
    Then the response status code should be 200
    And the updated todo should be marked as completed

  Scenario: Delete a todo
    Given the database contains a todo with title "Buy groceries"
    When I send a DELETE request for that todo
    Then the response status code should be 204
    And the todo should no longer exist in the database

  Scenario: Create a subtask
    Given the database contains a todo with title "Parent Task"
    When I send a POST request to add a subtask with title "Child Task" to "Parent Task"
    Then the response status code should be 201
    And the created subtask should have "Parent Task" as its parent
