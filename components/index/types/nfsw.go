package types

// NSFWClassification represents the classification returned by nsfw-server.
type NSFWClassification struct {
	Neutral float64 `json:"neutral"`
	Drawing float64 `json:"drawing"`
	Porn    float64 `json:"porn"`
	Hentai  float64 `json:"hentai"`
	Sexy    float64 `json:"sexy"`
}

// NSFW represents nsfw-server classification.
type NSFW struct {
	Classification    NSFWClassification `json:"classification"`
	NSFWServerVersion string             `json:"nsfwServerVersion"`
	ModelCID          string             `json:"modelCid"`
}
