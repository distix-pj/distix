package extractor

import (
	"fmt"
	"os"
	"errors"
	"strconv"

	rpm "github.com/sassoftware/go-rpmutils"
	"github.com/distix-pj/distix/data/model"
)

const S_ISDIR = 040000


type PkgExtractor struct {
	PkgPath	string
}

func NewPkgExtractor(pkgpath string) *PkgExtractor {
	return &PkgExtractor{
		PkgPath: pkgpath,
	}
}

func (e *PkgExtractor) Extract() (*model.Package, error) {
	fd, err := os.Open(e.PkgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	defer fd.Close()
	pkg, err := rpm.ReadRpm(fd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	pkgName, err := pkg.Header.GetString(rpm.NAME)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgVersion, err := pkg.Header.GetString(rpm.VERSION)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgRelease, err := pkg.Header.GetString(rpm.RELEASE)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	epochStr, err := pkg.Header.GetString(rpm.EPOCH)
	if err != nil {
		var noSuchTagErr rpm.NoSuchTagError
		if errors.As(err, &noSuchTagErr) {
			epochStr = ""
		} else {
			return nil, err
		}
	}
	var pkgEpoch *int
	if epochStr == "" {
		pkgEpoch = nil
	} else {
		val, err := strconv.Atoi(epochStr)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return nil, err
		}
		*pkgEpoch = val
	}
	pkgArch, err := pkg.Header.GetString(rpm.ARCH)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgSummary, err := pkg.Header.GetString(rpm.SUMMARY)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgDescription, err := pkg.Header.GetString(rpm.DESCRIPTION)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	rpmRequires, err := extractRequires(pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	rpmProvides, err := extractProvides(pkg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}

	filesInfo, err := pkg.Header.GetFiles()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	algo, err := pkg.Header.GetInt(rpm.FILEDIGESTALGO)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	digestAlgo := model.DigestAlgorithm(algo)
	provideFiles := []model.RpmFile{}
	for _, fileInfo := range filesInfo {
		if fileInfo.Mode() &^ 07777 == S_ISDIR { // Skip directories; only regular files are relevant for SBOM file entries.
			continue
		}
		provideFile := model.RpmFile{
			Name: fileInfo.Name(),
			Digest: fileInfo.Digest(),
			DigestAlgorithm: digestAlgo,
		}
		provideFiles = append(provideFiles, provideFile)
	}

	pkgSbomData := &model.Package{
		PkgNevra: model.PackageNevra{
			Name: pkgName,
			Epoch: pkgEpoch,
			Version: pkgVersion,
			Release: pkgRelease,
			Arch: pkgArch,
		},
		Summary: pkgSummary,
		Description: pkgDescription,
		Provides: rpmProvides,
		Requires: rpmRequires,
		Files: provideFiles,
	}
	return pkgSbomData, nil
}


func extractProvides(pkg *rpm.Rpm) ([]model.RpmCapability, error) {
	pkgProvideFlags, err := pkg.Header.GetUint32s(rpm.PROVIDEFLAGS)			/* 1112 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgProvideName, err := pkg.Header.GetStrings(rpm.PROVIDENAME)				/* 1047 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgProvideVersion, err := pkg.Header.GetStrings(rpm.PROVIDEVERSION) /* 1113 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	length := len(pkgProvideName)
	if len(pkgProvideFlags) != length || length != len(pkgProvideVersion) {
		return nil, fmt.Errorf("Array lengths don't match: flags=%d, names=%d, versions=%d",
			len(pkgProvideFlags), len(pkgProvideName), len(pkgProvideVersion))
	}
	provides := make([]model.RpmCapability, length)
	for i := 0; i < length; i++ {
		provides[i] = model.RpmCapability{
			Name:    pkgProvideName[i],
			Version: pkgProvideVersion[i],
			Flags:   pkgProvideFlags[i],
		}
	}
	return provides, nil
}


func extractRequires(pkg *rpm.Rpm) ([]model.RpmCapability, error) {
	pkgRequireFlags, err := pkg.Header.GetUint32s(rpm.REQUIREFLAGS)			/* 1048 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgRequireName, err := pkg.Header.GetStrings(rpm.REQUIRENAME)				/* 1049 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	pkgRequireVersion, err := pkg.Header.GetStrings(rpm.REQUIREVERSION) /* 1050 */
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, err
	}
	length := len(pkgRequireName)
	if len(pkgRequireFlags) != length || length != len(pkgRequireVersion) {
		return nil, fmt.Errorf("Array lengths don't match: flags=%d, names=%d, versions=%d",
			len(pkgRequireFlags), len(pkgRequireName), len(pkgRequireVersion))
	}

	requires := make([]model.RpmCapability, length)
	for i := 0; i < length; i++ {
		requires[i] = model.RpmCapability{
			Name:    pkgRequireName[i],
			Version: pkgRequireVersion[i],
			Flags:   pkgRequireFlags[i],
		}
	}
	return requires, nil
}

