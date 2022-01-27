package importer

// VisualisationUploaded provides an avro structure for a visualisation uploaded event
type VisualisationUploaded struct {
	Path string `avro:"path"`
}
