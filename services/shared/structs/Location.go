package structs

type Location struct {
	LocationType string    `json:"type" bson:"type"`
	LocationData []float32 `json:"coordinates" bson:"coordinates"`
}
