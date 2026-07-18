Feature: Wiki System
  As an authenticated project member
  I want to create and manage wiki pages
  So that I can document architecture and project details securely

  Background:
    Given I am a registered and authenticated user

  Scenario: Creating a new wiki page
    Given my account has a project named "Engineering"
    When I send an authenticated request to create a wiki page with title "Architecture" and slug "arch" and content "Base architecture"
    Then the wiki page should be saved successfully
    And I should see "Architecture" in the wiki index for "Engineering"

  Scenario: Prevent unauthorized access to wikis
    When I send an unauthenticated GET request to "/api/wikis"
    Then the response status code should be 401

  Scenario: Editing an existing wiki page
    Given my account has a wiki page with slug "arch"
    When I send an authenticated request to update the content to "Updated architecture"
    Then the wiki page should reflect the new content
    And the updated_at timestamp should change

  Scenario: Viewing a wiki page with markdown and shortcodes
    Given my account has a wiki page with content containing "{{diagram:1}}"
    And my account has a diagram with ID 1
    When I view the wiki page
    Then I should see the rendered markdown
    And the diagram 1 should be fetched and displayed seamlessly

  Scenario: Deleting a wiki page
    Given my account has a wiki page with slug "arch"
    When I send an authenticated request to delete the wiki page
    Then the wiki page should no longer appear in the wiki index
