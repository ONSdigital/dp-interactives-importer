package schema

import (
	"github.com/ONSdigital/dp-kafka/v3/avro"
)

var interactivesUploadedEvent = `{
  "type": "record",
  "name": "interactives-uploaded",
  "fields": [
	{"name": "collection_id", "type": "string"},
    {"name": "id", "type": "string"},
    {"name": "path", "type": "string"},
    {"name": "title", "type": "string"},
	{"name": "current_files", "type":["null",{"type":"array","items":"string"}]}
  ]
}`

// InteractivesUploadedEvent is the Avro schema for interactives uploaded messages.
var InteractivesUploadedEvent = &avro.Schema{
	Definition: interactivesUploadedEvent,
}
