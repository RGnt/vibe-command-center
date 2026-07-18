Feature: Mermaid Diagrams & Custom Icons
  As an application user
  I want to create diagrams and upload custom icons
  So that I can visualize workflows and architectures

  Scenario: Creating and saving a diagram
    Given I open the Diagram Editor
    When I type Mermaid syntax "graph TD; A-->B;"
    And I enter "My Flow" as the diagram name
    And I click save
    Then the diagram should be saved to the database
    And it should appear in the Load dropdown

  Scenario: Uploading a custom icon
    Given I am in the Diagram Editor
    When I select an image file to upload
    Then the image should be uploaded successfully
    And it should appear in the "Custom Icons" sidebar section under the "General" folder

  Scenario: Moving a custom icon to a folder
    Given a custom icon exists in the "General" folder
    When I click the move button and type "AWS"
    Then the icon should be moved to the "AWS" folder
    And the "AWS" folder should be visible and collapsible in the sidebar

  Scenario: Drag and drop custom icon to canvas
    Given I have an icon in the sidebar
    When I drag the icon to the canvas
    Then a new Mermaid node with the image tag should be added to the diagram
