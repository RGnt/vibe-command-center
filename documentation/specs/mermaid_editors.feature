Feature: Mermaid Diagram Editors
  As a user
  I want a split-pane diagram editor with code on the left and a live preview on the right
  So that I can write Mermaid source and immediately see the rendered result.

  Scenario: Split-pane layout is visible
    Given I am on the Mermaid Editor page
    Then the code editor panel should be visible on the left side
    And the live preview canvas should be visible on the right side

  Scenario: Applying code renders the diagram
    Given I am on the Mermaid Editor page
    When I type valid Mermaid code in the editor
    And I click "Apply Code"
    Then the live preview should update to show the rendered diagram

  Scenario: Zoom controls are available on the canvas
    Given I am on the Mermaid Editor page
    Then a zoom-in button should be visible on the canvas
    And a zoom-out button should be visible on the canvas
    And a reset zoom button should be visible on the canvas

  Scenario: Zooming in enlarges the diagram
    Given I am on the Mermaid Editor page
    When I click the zoom-in button multiple times
    Then the diagram should appear larger on the canvas

  Scenario: Resetting zoom returns to default scale
    Given I am on the Mermaid Editor page
    When I click the zoom-in button
    And I click the reset zoom button
    Then the diagram should return to its default scale

  Scenario: Dragging a connection between existing SVG nodes
    Given I am on the Mermaid Editor page
    And a diagram with at least two actors is rendered
    When I click an existing actor node
    Then a connection handle box appears next to the actor
    When I drag from the connection handle to another existing actor node
    Then a new sequence diagram connection edge syntax is appended to the Mermaid code

  Scenario: Saving a diagram
    Given I am on the Mermaid Editor page
    When I type a diagram name in the title input
    And I click "Save"
    Then the diagram should appear in the "Load Diagram" dropdown
