package multi

func makeMetadataProp(mimetype string) Properties {
	return map[string]interface{}{
		"metadata": map[string]interface{}{
			"Content-Type": mimetype,
		},
	}
}
