package schema

import (
	"github.com/ONSdigital/go-ns/avro"
)

var visualisationUploadedEvent = `{
  "type": "record",
  "name": "visualisation-uploaded",
  "fields": [
    {"name": "path", "type": "string", "default": ""},
  ]
}`

// ImageUploadedEvent is the Avro schema for Image uploaded messages.
var VisualisationUploadedEvent = &avro.Schema{
	Definition: visualisationUploadedEvent,
}
