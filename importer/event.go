package importer

// InteractivesUploaded provides an avro structure for an interactives uploaded event
type InteractivesUploaded struct {
	ID    string `avro:"id"`
	Path  string `avro:"path"`
	Title string `avro:"title"`
	CollectionID string `avro:"collection_id"`
}
