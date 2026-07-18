Feature: Wiki System
  As a project member
  I want to create and manage wiki pages
  So that I can document architecture and project details

  Scenario: Creating a new wiki page
    Given I have a project named "Engineering"
    When I create a wiki page with title "Architecture" and slug "arch" and content "Base architecture"
    Then the wiki page should be saved successfully
    And I should see "Architecture" in the wiki index for "Engineering"

  Scenario: Editing an existing wiki page
    Given a wiki page with slug "arch" exists
    When I update the content to "Updated architecture"
    Then the wiki page should reflect the new content
    And the updated_at timestamp should change

  Scenario: Viewing a wiki page with markdown and shortcodes
    Given a wiki page exists with content containing "{{diagram:1}}"
    And a diagram with ID 1 exists
    When I view the wiki page
    Then I should see the rendered markdown
    And the diagram 1 should be fetched and displayed seamlessly

  Scenario: Deleting a wiki page
    Given a wiki page with slug "arch" exists
    When I delete the wiki page
    Then the wiki page should no longer appear in the wiki index
