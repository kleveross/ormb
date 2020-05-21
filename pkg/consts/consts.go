package consts

const (
	ModelContentLayerMediaType = "application/tar+gzip"
	ModelConfigMediaType       = "application/vnd.caicloud.model.config.v1+json"

	ORMBfileName = "ormbfile.yaml"
)

// KnownMediaTypes returns a list of layer mediaTypes that the client knows about.
func KnownMediaTypes() []string {
	return []string{
		ModelContentLayerMediaType,
		ModelConfigMediaType,
	}
}
