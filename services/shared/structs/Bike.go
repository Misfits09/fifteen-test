package structs

type Bike struct {
	Id       string `json:"id" bson:"_id"`
	Location `json:"location" bson:"location"`
}
