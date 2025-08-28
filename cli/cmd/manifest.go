package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Haato3o/manifest/bytes"
	"github.com/Haato3o/manifest/manifest"
	"github.com/spf13/cobra"
	"os"
)

type createManifestArgs struct {
	fileName  string
	chunkSize bytes.Size
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Calculates and generates a manifest file",
	RunE: func(cmd *cobra.Command, args []string) error {
		chunkSize := cmd.Flag("chunk-size").Value.String()
		chunkSizeParsed, err := bytes.ParseSize(chunkSize)
		if err != nil {
			return err
		}

		cmdArgs := createManifestArgs{
			fileName:  args[0],
			chunkSize: chunkSizeParsed,
		}
		return create(cmdArgs)
	},
}

func create(args createManifestArgs) error {
	file, err := os.OpenFile(args.fileName, os.O_RDONLY, 0555)
	if err != nil {
		return err
	}
	defer file.Close()

	man, err := manifest.Create(args.fileName, args.chunkSize)
	if err != nil {
		return err
	}

	marshalled, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(marshalled))
	return nil
}

func init() {
	createCmd.Flags().StringP("chunk-size", "s", "1KB", "Chunk size in bytes")

	rootCmd.AddCommand(createCmd)
}
