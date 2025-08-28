package manifest

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/Haato3o/manifest/bytes"
	"math"
	"os"
)

var (
	emptyManifest = Manifest{}
	emptyFile     = File{}
)

func Create(folder string, chunkSize bytes.Size) (Manifest, error) {
	files, err := listFiles(folder)
	if err != nil {
		return emptyManifest, err
	}

	manifestFiles := make([]File, 0, len(files))
	for _, file := range files {
		manifestFile, err := createFile(file, chunkSize)
		if err != nil {
			return emptyManifest, err
		}

		manifestFiles = append(manifestFiles, manifestFile)
	}

	return Manifest{
		Files: manifestFiles,
	}, nil
}

func listFiles(folder string) ([]string, error) {
	files := make([]string, 0)

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryName := folder + "/" + entry.Name()
		if entry.IsDir() {
			children, err := listFiles(entryName)
			if err != nil {
				return nil, err
			}

			files = append(files, children...)
			continue
		}

		files = append(files, entryName)
	}

	return files, nil
}

func createFile(file string, chunkSize bytes.Size) (File, error) {
	reader, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err != nil {
		return emptyFile, err
	}
	defer reader.Close()

	stats, err := reader.Stat()
	if err != nil {
		return emptyFile, err
	}
	chunkByteSize := chunkSize.Bytes()
	chunksToRead := int(math.Ceil(float64(stats.Size()) / float64(chunkByteSize)))
	chunks := make([]Chunk, 0, chunksToRead)

	buffer := make([]byte, chunkByteSize)
	for i := 0; i < chunksToRead; i++ {
		if _, err = reader.Read(buffer); err != nil {
			return emptyFile, err
		}

		hash := createHash(buffer)

		chunks = append(chunks, Chunk{
			Hash: hash,
		})

		clear(buffer)
	}

	return File{
		Name:   file,
		Chunks: chunks,
	}, nil
}

func createHash(buffer []byte) string {
	hash := sha1.Sum(buffer)
	return hex.EncodeToString(hash[:])
}
