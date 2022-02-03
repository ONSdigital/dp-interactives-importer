package importer

// VisualisationUploaded provides an avro structure for a visualisation uploaded event
type VisualisationUploaded struct {
	ID   string `avro:"id"`
	Path string `avro:"path"`
}
