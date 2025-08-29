package manifest

type FileDiffResult struct {
	Name      string
	Different []Chunk
	Ratio     float64
}

type DiffResult struct {
	Different []FileDiffResult
	Ratio     float64
}

func Diff(leftManifest, rightManifest Manifest) DiffResult {
	differentFiles := diffFiles(leftManifest.Files, rightManifest.Files)

	actualRatio := 0.0
	totalRatio := 0.0
	for _, file := range differentFiles {
		actualRatio += file.Ratio
		totalRatio += 1.0
	}

	return DiffResult{
		Different: differentFiles,
		Ratio:     actualRatio / totalRatio,
	}
}

func diffFiles(leftFiles, rightFiles []File) []FileDiffResult {
	leftFilesMap := mapFiles(leftFiles)

	differentFiles := make([]FileDiffResult, 0)
	for _, rightFile := range rightFiles {
		leftFile, exists := leftFilesMap[rightFile.Name]
		if !exists {
			differentFiles = append(differentFiles, FileDiffResult{
				Name:      rightFile.Name,
				Different: rightFile.Chunks,
				Ratio:     1.0,
			})
			continue
		}

		differentChunks := diffChunks(leftFile.Chunks, rightFile.Chunks)
		differentFiles = append(differentFiles, FileDiffResult{
			Name:      rightFile.Name,
			Different: differentChunks,
			Ratio:     float64(len(differentChunks)) / float64(len(rightFile.Chunks)),
		})
	}

	return differentFiles
}

func diffChunks(leftChunks, rightChunks []Chunk) []Chunk {
	differentChunks := make([]Chunk, 0)
	// leftChunksIndexed := indexChunks(leftChunks)

	for idx, rightChunk := range rightChunks {
		if idx >= len(leftChunks) {
			differentChunks = append(differentChunks, rightChunk)
			continue
		}

		leftChunk := leftChunks[idx]
		if leftChunk.Hash != rightChunk.Hash {
			differentChunks = append(differentChunks, leftChunk)
		}
	}

	return differentChunks
}

func indexChunks(chunks []Chunk) map[string]Chunk {
	m := make(map[string]Chunk, len(chunks))

	for _, chunk := range chunks {
		m[chunk.Hash] = chunk
	}

	return m
}

func mapFiles(files []File) map[string]File {
	m := make(map[string]File, len(files))
	for _, file := range files {
		m[file.Name] = file
	}

	return m
}
