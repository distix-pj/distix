package command

import (
	"fmt"
	"os"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/distix-pj/distix/data/extractor"
	"github.com/distix-pj/distix/format"
)


type PackageRunner struct {
	InputFile string
}
var inputFile string

func (r *PackageRunner) Setup() error {
	absInputFilePath, err := filepath.Abs(inputFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	r.InputFile = absInputFilePath
	slog.Debug("PackageRunner Options: ",
		"options", r,
		"InputFile", inputFile,
	)
	return nil
}

func (r *PackageRunner) Run() error {
	ext := extractor.NewPkgExtractor(r.InputFile)
	pkgData, err := ext.Extract()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	slog.Debug("pkgData: %s", pkgData)

	w, err := format.NewWriter(RootOpts.SbomType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return w.WritePackage(pkgData, RootOpts.OutputFile)
}


func NewPackageCmd() *cobra.Command {
	runner := &PackageRunner{}
	cmd := &cobra.Command{
		Use:   "package",
		Short: "Generate SBOM from RPM Package",
		// Long: `A longer description here..`
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runner.Setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Run()
		},
	}
	cmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "Path to target RPM Package")
	return cmd
}

func init() {
	RegisterSubCommand(NewPackageCmd())
}

