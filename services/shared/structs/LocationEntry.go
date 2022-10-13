package structs

type InternalLocationEntry struct {
	Id       string `json:"id" bson:"bikeId"`
	Location `json:"location" bson:"location"`
	Time     int64 `bson:"time" json:"time"`
}

type APILocationEntry struct {
	Id        string `json:"id" bson:"bikeId"`
	Location  `json:"location" bson:"location"`
	TimeStamp string `bson:"time" json:"time"`
}
