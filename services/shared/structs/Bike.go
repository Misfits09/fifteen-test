package structs

type Bike struct {
	ID       string `json:"id" bson:"_id"`
	Location `json:"location" bson:"location"`
}
