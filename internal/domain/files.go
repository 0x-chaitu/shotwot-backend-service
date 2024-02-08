package domain

type File struct {
	Name     string `bson:"name" json:"name"`
	Filetype string `bson:"type" json:"type"`
}
