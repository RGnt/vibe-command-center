Feature: User Authentication
  In order to secure their workspace and personal data
  As a user
  I want to register an account and log in using my email and password

  Scenario: Register a new account
    Given I am a new user
    When I submit my email and a secure password
    Then a new account should be created for me
    And I should be automatically logged in
    And a default project, workflow, and wiki page should be generated

  Scenario: Login with valid credentials
    Given I have an existing account
    When I submit my registered email and password
    Then I should receive a JWT authentication token
    And I should be redirected to my workspace

  Scenario: Login with invalid credentials
    Given I have an existing account
    When I submit an incorrect password
    Then I should see an authentication error
    And I should not be granted access

  Scenario: Prevent duplicate email registration
    Given an account exists with the email "user@example.com"
    When I attempt to register a new account with "user@example.com"
    Then I should see an error indicating the email might already exist
