package manifest

type Chunk struct {
	ID   int    `json:"id"`
	Hash string `json:"hash"`
}

type File struct {
	Name   string  `json:"name"`
	Chunks []Chunk `json:"chunks"`
}

type Manifest struct {
	Files []File `json:"files"`
}
