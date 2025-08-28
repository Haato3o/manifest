package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Haato3o/manifest/manifest"
	"github.com/spf13/cobra"
	"io"
	"os"
)

type diffArgs struct {
	left  string
	right string
}

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compares two manifests and returns the difference between them",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmdArgs := diffArgs{
			left:  cmd.Flag("left").Value.String(),
			right: cmd.Flag("right").Value.String(),
		}

		return diff(cmdArgs)
	},
}

func diff(args diffArgs) error {
	leftManifest, err := loadManifest(args.left)
	if err != nil {
		return err
	}

	rightManifest, err := loadManifest(args.right)
	if err != nil {
		return err
	}

	result := manifest.Diff(leftManifest, rightManifest)

	for _, file := range result.Different {
		fmt.Printf("File name: %s\n", file.Name)
		fmt.Printf("Equality Ratio: %.2f\n", (1-file.Ratio)*100)
		fmt.Println("_____________________________________")
	}
	fmt.Printf("Total equality ratio: %.2f\n", (1-result.Ratio)*100)

	return nil
}

func loadManifest(path string) (manifest.Manifest, error) {
	man := manifest.Manifest{}

	file, err := os.OpenFile(path, os.O_RDONLY, 0555)
	if err != nil {
		return man, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return man, err
	}

	if err = json.Unmarshal(fileContent, &man); err != nil {
		return man, err
	}

	return man, nil
}

func init() {
	diffCmd.Flags().StringP("left", "l", "", "Left manifest")
	diffCmd.Flags().StringP("right", "r", "", "Right manifest")

	_ = diffCmd.MarkFlagRequired("left")
	_ = diffCmd.MarkFlagRequired("right")

	rootCmd.AddCommand(diffCmd)
}
