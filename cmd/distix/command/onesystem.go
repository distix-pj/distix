package command

import (
	"fmt"
	"os"
	"log/slog"
	"path/filepath"

	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"

	"github.com/distix-pj/distix/data/extractor"
	"github.com/distix-pj/distix/format"
)

const DEF_RPMDB_PATH = "/var/lib/rpm/rpmdb.sqlite"


type OneSystemRunner struct {
	RpmDb string
}
var oneSysRpmDb string

func (r *OneSystemRunner) Setup() error {
	absRpmDbPath, err := filepath.Abs(oneSysRpmDb)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	r.RpmDb = absRpmDbPath

	slog.Debug("OneSystemRunner Options: ",
		"options", r,
		"RpmDb", oneSysRpmDb,
	)
	return nil
}

func (r *OneSystemRunner) Run() error {
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
	return w.WriteOneSystem(sysData, RootOpts.OutputFile)
}


func NewOneSystemCmd() *cobra.Command {
	runner := &OneSystemRunner{}
	cmd := &cobra.Command{
		Use:   "onesystem",
		Short: "Generate SBOM for system in one SBOM file (from RPMDB)",
		// Long: `A longer description here..`
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return runner.Setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Run()
		},
	}
	cmd.Flags().StringVarP(&oneSysRpmDb, "rpmdb", "r", DEF_RPMDB_PATH, "Path to target RPM Package")
	return cmd
}

func init() {
	RegisterSubCommand(NewOneSystemCmd())
}

