package db

type CustomizablePrice struct {
	Value        string            `bson:"value"`
	CustomValues map[string]string `bson:"customValues" json:"customValues"`
}
