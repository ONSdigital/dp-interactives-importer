package importer

// InteractivesUploaded provides an avro structure for a interactives uploaded event
type InteractivesUploaded struct {
	ID           string   `avro:"id"`
	Path         string   `avro:"path"`
	Title        string   `avro:"title"`
	CurrentFiles []string `avro:"current_files"`
}
