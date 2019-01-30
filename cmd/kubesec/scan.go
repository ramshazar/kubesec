package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/sublimino/kubesec/pkg/ruler"
	"io/ioutil"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:     `scan [file]`,
	Short:   "Scans Kubernetes resource YAML or JSON",
	Example: `  scan ./deployment.yaml`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file path is required")
		}

		filename, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		fileBytes, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var data []byte
		isJson := json.Valid(fileBytes)
		if isJson {
			data = fileBytes
		} else {
			data, err = yaml.YAMLToJSON(fileBytes)
			if err != nil {
				return err
			}
		}

		report := ruler.NewRuleset(logger).Run(data)
		res, err := json.Marshal(report)
		if err != nil {
			return err
		}

		fmt.Println(prettyJSON(res))
		return nil
	},
}

func prettyJSON(b []byte) string {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return string(out.Bytes())
}
