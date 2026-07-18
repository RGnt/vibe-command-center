Feature: Todo Management
  As an authenticated user
  I want to be able to manage my todos
  So that I can keep track of my tasks securely and separate from others

  Background:
    Given I am a registered and authenticated user

  Scenario: Create a new todo
    Given the backend server is running
    When I send an authenticated POST request to "/api/todos" with the following JSON:
      """
      {
        "title": "Buy groceries",
        "content": "Milk, eggs, bread",
        "stage": "To Do"
      }
      """
    Then the response status code should be 201
    And the response should contain a todo with title "Buy groceries", content "Milk, eggs, bread", and stage "To Do"
    And the todo should belong to me

  Scenario: Retrieve all my todos
    Given my account contains a todo with title "Buy groceries"
    When I send an authenticated GET request to "/api/todos"
    Then the response status code should be 200
    And the response should contain a list of todos
    And the list should include a todo with title "Buy groceries"

  Scenario: Prevent unauthorized access to todos
    When I send an unauthenticated GET request to "/api/todos"
    Then the response status code should be 401

  Scenario: Update my todo
    Given my account contains a todo with title "Buy groceries" and stage "To Do"
    When I send an authenticated PUT request to update the todo to:
      """
      {
        "title": "Buy groceries",
        "stage": "In Progress"
      }
      """
    Then the response status code should be 200
    And the updated todo should have the stage "In Progress"

  Scenario: Toggle my todo completion status
    Given my account contains an incomplete todo with title "Buy groceries"
    When I send an authenticated PATCH request to toggle the todo completion
    Then the response status code should be 200
    And the updated todo should be marked as completed

  Scenario: Delete my todo
    Given my account contains a todo with title "Buy groceries"
    When I send an authenticated DELETE request for that todo
    Then the response status code should be 204
    And the todo should no longer exist in the database

  Scenario: Create a subtask
    Given my account contains a todo with title "Parent Task"
    When I send an authenticated POST request to add a subtask with title "Child Task" to "Parent Task"
    Then the response status code should be 201
    And the created subtask should have "Parent Task" as its parent
