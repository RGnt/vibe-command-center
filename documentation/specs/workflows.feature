Feature: Workflow Management
  As an authenticated user
  I want to be able to manage custom workflows
  So that I can tailor the Kanban board columns to my needs

  Background:
    Given I am a registered and authenticated user

  Scenario: Retrieve all workflows
    Given my account contains a default workflow
    When I send an authenticated GET request to "/api/workflows"
    Then the response status code should be 200
    And the response should contain a list of workflows including "Default Workflow"

  Scenario: Prevent unauthorized access to workflows
    When I send an unauthenticated GET request to "/api/workflows"
    Then the response status code should be 401

  Scenario: Create a new custom workflow
    Given the backend server is running
    When I send an authenticated POST request to "/api/workflows" with the following JSON:
      """
      {
        "name": "Software Development",
        "stages": [
          {"name": "Backlog"},
          {"name": "In Dev"},
          {"name": "QA"},
          {"name": "Deployed"}
        ]
      }
      """
    Then the response status code should be 201
    And the response should contain a workflow with name "Software Development"
    And it should have 4 stages in the correct order

  Scenario: Update a workflow
    Given my account contains a workflow with name "Custom Workflow"
    When I send an authenticated PUT request to update the workflow to have stages "Todo", "Doing", "Done"
    Then the response status code should be 200
    And the workflow should be updated to contain those 3 stages

  Scenario: Delete a workflow
    Given my account contains a workflow with name "Legacy Workflow"
    When I send an authenticated DELETE request for that workflow
    Then the response status code should be 204
    And the workflow should no longer exist in the database
