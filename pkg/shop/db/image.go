package db

type Image struct {
	URL     string `bson:"url" json:"url"`
	Title   string `bson:"title" json:"title"`
	AltText string `bson:"altText" json:"altText"`
}
