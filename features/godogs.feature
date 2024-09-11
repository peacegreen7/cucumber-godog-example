Feature: Register user

  Scenario: Register user successfully
    Given I send a "POST" request to "https://reqres.in/api/register" with payload:
    """
    {
      "email": "eve.holt@reqres.in",
      "password": "pistol"
    }
    """
    Then the response code should be 200
    And the response should match json:
    """
    {
      "id":4,
      "token":"QpwL5tke4Pnpja7X4"
    }
    """
