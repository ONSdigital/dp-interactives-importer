Feature: Consuming messages from Kafka

  Scenario: dp-interactives-api has sent one valid message
    Given these events are consumed:
      | ID      | Path                |
      | valid-1 | /some/path/name.css |
    Then "1" interactives should be downloaded from s3 successfully
    And "1" interactives should be uploaded via the upload service