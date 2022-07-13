package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var interactivesUploadedEvent = `{
  "type": "record",
  "name": "interactives-uploaded",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "path", "type": "string"},
    {"name": "title", "type": "string"}
  ]
}`

// InteractivesUploadedEvent is the Avro schema for interactives uploaded messages.
var InteractivesUploadedEvent = &avro.Schema{
	Definition: interactivesUploadedEvent,
}
