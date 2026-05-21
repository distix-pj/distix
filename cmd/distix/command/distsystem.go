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

type DistSystemRunner struct {
	RpmDb string
	OutputSubDir string
}
var distSysRpmDb string
var outputSubDir string

func (r *DistSystemRunner) Setup() error {
	absRpmDbPath, err := filepath.Abs(distSysRpmDb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	r.RpmDb = absRpmDbPath

	absOutputSubDir, err := filepath.Abs(outputSubDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absOutputSubDir, 0755); err != nil {
		return err
	}
	r.OutputSubDir = absOutputSubDir

	slog.Debug("PackageRunner Options: ",
		"options", r,
		"RpmDb", distSysRpmDb,
		"outputSubDir", outputSubDir,
	)
	return nil
}

func (r *DistSystemRunner) Run() error {
	ext := extractor.NewRpmdbExtractor(r.RpmDb)
	sysData, err := ext.Extract()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	w, err := format.NewWriter(RootOpts.SbomType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	if err := w.WriteDistSystem(sysData, RootOpts.OutputFile, r.OutputSubDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	for _, pkg := range sysData.Packages {
		fileName := fmt.Sprintf(
			"%s.%s.%s",
			pkg.PkgNevra.GetNEVRA(),
			RootOpts.SbomType.RecordType,
			RootOpts.SbomType.FileFormatType,
		)
		outputPath := filepath.Join(r.OutputSubDir, fileName)
		fd, err := os.Create(outputPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
		if err := w.WritePackage(&pkg, fd); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
	}

	return nil
}


func NewDistSystemCmd() *cobra.Command {
	runner := &DistSystemRunner{}
	cmd := &cobra.Command{
		Use:   "distsystem",
		Short: "Generate SBOM for system using external references to each package SBOM.",
		// Long: `A longer description here..`
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runner.Setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Run()
		},
	}
	cmd.Flags().StringVarP(&distSysRpmDb, "rpmdb", "r", DEF_RPMDB_PATH, "Path to target RPM Package")
	cmd.Flags().StringVarP(&outputSubDir, "output-subdir", "O", "pkg-sboms", "Output Sub Dir (required)")
	return cmd
}

func init() {
	RegisterSubCommand(NewDistSystemCmd())
}
