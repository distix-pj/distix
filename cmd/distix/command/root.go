package command

import (
	"io"
	"os"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/distix-pj/distix/format/types"
)

var rootCmd *cobra.Command


type RootRunner struct {
  Verbose bool
	SbomType types.SbomType
	OutputFile io.Writer
	OutputSubDir string
}
var RootOpts *RootRunner
var outputFile string = ""
var outputSubDir string

func (r *RootRunner) Setup() error {
	loglevel := slog.LevelInfo
	if r.Verbose {
		loglevel = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: loglevel,
	})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	
	if outputFile == "" {
		r.OutputFile = os.Stdout
	} else {
		absPath, err := filepath.Abs(outputFile)
		if err != nil {
			return err
		}
		file, err := os.Create(absPath)
		if err != nil {
			return err
		}
		r.OutputFile = file
	}

	absOutputSubDir, err := filepath.Abs(outputSubDir)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(absOutputSubDir, 0755); err != nil {
		return err
	}
	r.OutputSubDir = absOutputSubDir

	slog.Debug("RootRunner Options",
		"options", r,
		"verbose", r.Verbose,
		"sbomType", r.SbomType,
		"outputFile", outputFile,
		"outputSubDir", outputSubDir,
	)

	return nil
}


func NewRootCmd() *cobra.Command {
	runner := &RootRunner{
		// Verbose: false,
		SbomType: types.GetSbomTypeDefault(),
	}
	RootOpts = runner
	cmd := &cobra.Command{
		Use:   "distix",
		Short: "distix",
		// Long: `A longer description here..`
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return runner.Setup()
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			if closer, ok := RootOpts.OutputFile.(io.Closer); ok {
				return closer.Close()
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolVarP(&runner.Verbose, "verbose", "v", false, "Verbose message")
	cmd.PersistentFlags().VarP(&runner.SbomType, "format-type", "", "Sbom Type")
	cmd.PersistentFlags().StringVarP(&outputFile, "output-file", "o", "", "Output file (default stdout)")
	cmd.PersistentFlags().StringVarP(&outputSubDir, "output-subdir", "O", "subcomps", "Output Sub Dir (required)")
	return cmd
}

func Execute() error {
	err := rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}

func initRootCmd() {
	if rootCmd == nil {
		rootCmd = NewRootCmd()
	}
}

func RegisterSubCommand(subCmd *cobra.Command) {
	initRootCmd()
	rootCmd.AddCommand(subCmd)
}

func init() {
	initRootCmd()
}

