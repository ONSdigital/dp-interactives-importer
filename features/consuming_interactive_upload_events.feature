Feature: Consuming messages from Kafka

  Scenario: dp-interactives-api has sent one valid message
    Given these events are consumed:
      | ID      | Path            |
      | valid-1 | /some/path/name |
    Then "1" interactives should be uploaded to s3 successfully
    And a message for "valid-1" with "s3:///some/path/name" is produced