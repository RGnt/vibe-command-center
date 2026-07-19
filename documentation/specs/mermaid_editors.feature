Feature: Mermaid Diagram Editors
  As a user
  I want to have a floating toolbar with a shape palette and a collapsible code editor for sequence diagrams
  So that I can visually create sequence diagrams efficiently on a full-width canvas.

  Scenario: Shape palette displays sequence-specific tools
    Given I am on the Mermaid Editor page
    When I click "Add Shape" on the floating toolbar
    Then a modal should display sequence-specific tools (Actor, Participant, Note)

  Scenario: Clicking a shape adds it to the canvas
    Given I am on the Mermaid Editor page
    When I click "Add Shape" on the floating toolbar
    And I select the "Actor" shape from the modal
    Then a new "Actor" node should appear on the canvas

  Scenario: Collapsible code editor
    Given I am on the Mermaid Editor page
    When I click the "Toggle Code Editor" button on the left edge of the canvas
    Then the code preview panel should expand
    When I click the "Toggle Code Editor" button again
    Then the code preview panel should collapse

  Scenario: Dragging a connection between existing SVG nodes
    Given I am on the Mermaid Editor page
    When I click an existing actor node
    Then a connection handle box appears next to the actor
    When I drag from the connection handle to another existing actor node
    Then a new sequence diagram connection edge syntax is appended to the Mermaid code

  Scenario: Dragging a connection to empty space creates a new actor
    Given I am on the Mermaid Editor page
    When I click an existing actor node
    And I drag from the connection handle to an empty space on the canvas
    Then a shape palette modal appears under my cursor
    When I select the "Actor" shape from the palette
    Then the new actor and connection syntax is appended to the Mermaid code
