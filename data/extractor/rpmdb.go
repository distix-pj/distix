package extractor

import (
	"fmt"
	"os"

	rpmdb "github.com/knqyf263/go-rpmdb/pkg"

	"github.com/distix-pj/distix/data/model"
)


type RpmdbExtractor struct {
	RpmdbPath	string
}

func NewRpmdbExtractor(rpmdbpath string) *RpmdbExtractor {
	return &RpmdbExtractor{
		RpmdbPath: rpmdbpath,
	}
}

func (e *RpmdbExtractor) Extract() (*model.System, error) {
	db, err := rpmdb.Open(e.RpmdbPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	defer db.Close()

	pkgList, err := db.ListPackages()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgs := []model.Package{}

	for _, pkg := range pkgList {
		if pkg.Name == "gpg-pubkey" {
			continue
		}

		rpmRequires, err := extractRequiresFromPackageInfo(pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		rpmProvides, err := extractProvidesFromPackageInfo(pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		rpmFiles, err := extractFilesFromPackageInfo(pkg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}

		pkgSbomData := &model.Package{
			PkgNevra: model.PackageNevra{
				Name: pkg.Name,
				// Epoch: pkg.EpochNum(),
				Epoch: pkg.Epoch,
				Version: pkg.Version,
				Release: pkg.Release,
				Arch: pkg.Arch,
			},
			// Summary: pkgSummary,
			// Description: pkgDescription,
			Provides: rpmProvides,
			Requires: rpmRequires,
			Files: rpmFiles,
		}
		pkgs = append(pkgs, *pkgSbomData)
	}

	return &model.System{
		HostName: "test hostname",
		Packages: pkgs,
	}, nil
}


func extractRequiresFromPackageInfo(pkg *rpmdb.PackageInfo) ([]model.RpmCapability, error) {
	requires := make([]model.RpmCapability, len(pkg.Requires))
	for i := 0; i < len(pkg.Requires); i++ {
		requires[i] = model.RpmCapability{
			Name:    pkg.Requires[i],
			// Version: pkgRequireVersion[i],
			// Flags:   pkgRequireFlags[i],
		}
	}
	return requires, nil
}

func extractProvidesFromPackageInfo(pkg *rpmdb.PackageInfo) ([]model.RpmCapability, error) {
	provides := make([]model.RpmCapability, len(pkg.Provides))
	for i := 0; i < len(pkg.Provides); i++ {
		provides[i] = model.RpmCapability{
			Name:    pkg.Provides[i],
			// Version: pkgRequireVersion[i],
			// Flags:   pkgRequireFlags[i],
		}
	}
	return provides, nil
}

func extractFilesFromPackageInfo(pkg *rpmdb.PackageInfo) ([]model.RpmFile, error) {
	ffs, err := pkg.InstalledFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	provideFiles := make([]model.RpmFile, len(ffs))
	for i, file := range ffs {
		provideFiles[i] = model.RpmFile{
			Name: file.Path,
			Digest: file.Digest,
			DigestAlgorithm: model.DigestAlgorithm(pkg.DigestAlgorithm),
		}
	}
	return provideFiles, nil
}

