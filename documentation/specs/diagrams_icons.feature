Feature: Diagrams and Icons
  As an authenticated user
  I want to create diagrams and upload custom icons
  So that I can use them in my wikis and workflows

  Background:
    Given I am a registered and authenticated user

  Scenario: Creating a mermaid diagram
    When I send an authenticated request to create a new diagram with type "graph TD" and some mermaid code
    Then the diagram should be saved to my account successfully
    And I can reference it by ID in my wikis

  Scenario: Editing a diagram
    Given my account has a diagram with ID "uuid"
    When I send an authenticated request to update the diagram code
    Then the diagram should be updated
    And any wiki referencing it will show the new diagram

  Scenario: Prevent unauthorized access to diagrams
    When I send an unauthenticated GET request to "/api/diagrams"
    Then the response status code should be 401

  Scenario: Uploading a custom icon
    When I send an authenticated POST request with a file "icon.png" to "/api/upload"
    Then the file should be saved in my account
    And I should receive a URL to the uploaded file
    And a new Icon record should be created in the database for me

  Scenario: Deleting an icon
    Given my account has an icon with ID "uuid"
    When I send an authenticated request to delete the icon
    Then the icon record should be removed
    And it should no longer be available in the icon picker
