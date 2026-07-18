Feature: Project Management
  In order to organize different areas of work
  As a user
  I want to create and manage multiple projects

  Scenario: Create a new project
    Given I am logged in
    When I create a new project with a name, description, and workflow ID
    Then the project should be saved to my workspace
    And I should be able to navigate to the new project board

  Scenario: View my projects
    Given I have created multiple projects
    When I view my projects list
    Then I should see all of the projects I have access to

  Scenario: Update an existing project
    Given I have an existing project
    When I update its name or description
    Then the changes should be reflected immediately

  Scenario: Delete a project
    Given I have an existing project
    When I delete the project
    Then the project and all its associated tasks and wikis should be removed
