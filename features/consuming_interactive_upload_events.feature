Feature: Consuming messages from Kafka

  Scenario: dp-interactives-api has sent a message with an invalid interactive id
    Given these events are consumed:
      | ID        | Path                         |
      | invalid-1 | test_zips/does_not_exist.zip |
    Then "0" interactives should be uploaded via the upload service
    And "invalid-1" interactive should be updated as a failure via the interactives API

  Scenario: dp-interactives-api has sent a message with a non-existent zip file
    Given these events are consumed:
      | ID      | Path                     |
      | valid-1 | test_zips/does_not_exist.zip |
    Then "1" interactives should be downloaded from s3 successfully
    And "0" interactives should be uploaded via the upload service
    And "valid-1" interactive should be updated as a failure via the interactives API

  Scenario: dp-interactives-api has sent a message with a valid zip file
    Given these events are consumed:
      | ID      | Path                     |
      | valid-1 | test_zips/happy_path.zip |
    Then "1" interactives should be downloaded from s3 successfully
    And "11" interactives should be uploaded via the upload service
    And "valid-1" interactive should be successfully updated via the interactives API

  Scenario: dp-interactives-api has sent a message with an empty zip file
    Given these events are consumed:
      | ID      | Path                     |
      | valid-1 | test_zips/empty.zip |
    Then "1" interactives should be downloaded from s3 successfully
    And "0" interactives should be uploaded via the upload service
    And "valid-1" interactive should be updated as a failure via the interactives API

  Scenario: dp-interactives-api has sent a message with an corrupt zip file
    Given these events are consumed:
      | ID      | Path                     |
      | valid-1 | test_zips/random_bytes.zip |
    Then "1" interactives should be downloaded from s3 successfully
    And "0" interactives should be uploaded via the upload service
    And "valid-1" interactive should be updated as a failure via the interactives API

  Scenario: dp-interactives-api has sent a message with an invalid zip file
    Given these events are consumed:
      | ID      | Path                     |
      | valid-1 | test_zips/bad_content.zip |
    Then "1" interactives should be downloaded from s3 successfully
    And "0" interactives should be uploaded via the upload service
    And "valid-1" interactive should be updated as a failure via the interactives API