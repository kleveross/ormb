package consts

const (
	// MediaTypeModelContentLayer is the media type of the content layer.
	MediaTypeModelContentLayer = "application/tar+gzip"
	// MediaTypeModelConfig is the media type of the model config.
	MediaTypeModelConfig = "application/vnd.caicloud.model.config.v1alpha1+json"

	// ORMBfileName is the filename of the config file.
	ORMBfileName = "ormbfile.yaml"

	// ORMBModelDirectory is the directory name of the model in the path.
	ORMBModelDirectory = "model"
)

// KnownMediaTypes returns a list of layer mediaTypes that the client knows about.
func KnownMediaTypes() []string {
	return []string{
		MediaTypeModelConfig,
		MediaTypeModelContentLayer,
	}
}
