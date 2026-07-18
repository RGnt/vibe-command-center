Feature: Workflow Management
  As a user
  I want to be able to manage custom workflows
  So that I can tailor the Kanban board columns to my needs

  Scenario: Retrieve all workflows
    Given the database contains a default workflow
    When I send a GET request to "/api/workflows"
    Then the response status code should be 200
    And the response should contain a list of workflows including "Default Workflow"

  Scenario: Create a new custom workflow
    Given the backend server is running
    When I send a POST request to "/api/workflows" with the following JSON:
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
    Given the database contains a workflow with name "Custom Workflow"
    When I send a PUT request to update the workflow to have stages "Todo", "Doing", "Done"
    Then the response status code should be 200
    And the workflow should be updated to contain those 3 stages

  Scenario: Delete a workflow
    Given the database contains a workflow with name "Legacy Workflow"
    When I send a DELETE request for that workflow
    Then the response status code should be 204
    And the workflow should no longer exist in the database
