package format

import (
	"errors"
	"io"
	"fmt"
	"time"

	spdxv2_3 "github.com/spdx/tools-golang/spdx/v2/v2_3"
	spdxcommon "github.com/spdx/tools-golang/spdx/v2/common"
	spdxjson "github.com/spdx/tools-golang/json"

	"github.com/distix-pj/distix/data/model"
)

type SPDXWriter struct {
	fileFormat SbomFileFormatType
}

func (w *SPDXWriter) WritePackage(pkg *model.Package, out io.Writer) error {
	pkgElemID := spdxcommon.ElementID("SPDXRef-Package")

	packages := []*spdxv2_3.Package{
		{
			PackageName:             pkg.PkgNevra.Name,
			PackageVersion:          pkg.PkgNevra.Version,
			PackageSPDXIdentifier:   pkgElemID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
			PackageExternalReferences: []*spdxv2_3.PackageExternalReference{
				{
					Category: "PACKAGE-MANAGER",
					RefType:  "purl",
					Locator:  pkg.GetPurl(),
				},
			},
		},
	}
	relationships := []*spdxv2_3.Relationship{
		{
			RefA:         spdxcommon.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         spdxcommon.DocElementID{ElementRefID: pkgElemID},
			Relationship: "DESCRIBES",
		},
	}

	for i, prov := range pkg.Provides {
		provID := spdxcommon.ElementID(fmt.Sprintf("SPDXRef-Prov-%d", i))
		packages = append(packages, &spdxv2_3.Package{
			PackageName:             prov.Name,
			PackageVersion:          prov.Version,
			PackageSPDXIdentifier:   provID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
		})
		relationships = append(relationships, &spdxv2_3.Relationship{
			RefA:         spdxcommon.DocElementID{ElementRefID: pkgElemID},
			RefB:         spdxcommon.DocElementID{ElementRefID: provID},
			Relationship: "CONTAINS",
		})
	}

	for i, req := range pkg.Requires {
		reqID := spdxcommon.ElementID(fmt.Sprintf("SPDXRef-Req-%d", i))
		packages = append(packages, &spdxv2_3.Package{
			PackageName:             req.Name,
			PackageVersion:          req.Version,
			PackageSPDXIdentifier:   reqID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
		})
		relationships = append(relationships, &spdxv2_3.Relationship{
			RefA:         spdxcommon.DocElementID{ElementRefID: pkgElemID},
			RefB:         spdxcommon.DocElementID{ElementRefID: reqID},
			Relationship: "DEPENDS_ON",
		})
	}

	files := []*spdxv2_3.File{}
	for i, f := range pkg.Files {
		fileID := spdxcommon.ElementID(fmt.Sprintf("SPDXRef-File-%d", i))
		algo, err := toSPDXAlgorithm(f.DigestAlgorithm)
		if err != nil {
			return err
		}

		files = append(files, &spdxv2_3.File{
			FileName: f.Name,
			FileSPDXIdentifier: fileID,
			Checksums: []spdxcommon.Checksum{{
					Algorithm: algo,
					Value: f.Digest,
			}},
		})
		relationships = append(relationships, &spdxv2_3.Relationship{
			RefA:         spdxcommon.DocElementID{ElementRefID: pkgElemID},
			RefB:         spdxcommon.DocElementID{ElementRefID: fileID},
			Relationship: "CONTAINS",
		})
	}

	doc := &spdxv2_3.Document{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXIdentifier:    "DOCUMENT",
		DocumentName:      pkg.PkgNevra.GetNEVRA(),
		DocumentNamespace: "https://distix.example.org/package/" + pkg.GetPurl(),
		CreationInfo: &spdxv2_3.CreationInfo{
			Created: time.Now().UTC().Format(time.RFC3339),
			Creators: []spdxcommon.Creator{
				{CreatorType: "Tool", Creator: "distix"},
			},
		},
		Packages:          packages,
		Relationships:     relationships,
		Files:             files,
	}

	return w.write(doc, out)
}

func (w *SPDXWriter) WriteOneSystem(sys *model.System, out io.Writer) error {
	sysID := spdxcommon.ElementID("SPDXRef-System")
	rpmPkgsID := spdxcommon.ElementID("SPDXRef-RPMPackages")
	packages := []*spdxv2_3.Package{
		{
			PackageName:             sys.HostName,
			PackageSPDXIdentifier:   sysID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
		},
		{
			PackageName:             "RPM-Packages",
			PackageSPDXIdentifier:   rpmPkgsID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
		},
	}
	relationships := []*spdxv2_3.Relationship{
		{
			RefA:         spdxcommon.DocElementID{ElementRefID: "DOCUMENT"},
			RefB:         spdxcommon.DocElementID{ElementRefID: sysID},
			Relationship: "DESCRIBES",
		},
		{
			RefA:         spdxcommon.DocElementID{ElementRefID: sysID},
			RefB:         spdxcommon.DocElementID{ElementRefID: rpmPkgsID},
			Relationship: "CONTAINS",
		},
	}

	providerMap := map[string]spdxcommon.ElementID{}
	for _, pkg := range sys.Packages {
		pkgID := spdxcommon.ElementID(fmt.Sprintf("SPDXRef-Pkg-%s", pkg.GetPurl()))
		for _, prov := range pkg.Provides {
			providerMap[prov.Name] = pkgID
		}

		packages = append(packages, &spdxv2_3.Package{
			PackageName:             pkg.PkgNevra.Name,
			PackageVersion:          pkg.PkgNevra.Version,
			PackageSPDXIdentifier:   pkgID,
			PackageDownloadLocation: "NOASSERTION",
			FilesAnalyzed:           false,
			PackageExternalReferences: []*spdxv2_3.PackageExternalReference{
				{
					Category: "PACKAGE-MANAGER",
					RefType:  "purl",
					Locator:  pkg.GetPurl(),
				},
			},
		})
		relationships = append(relationships, &spdxv2_3.Relationship{
			RefA:         spdxcommon.DocElementID{ElementRefID: rpmPkgsID},
			RefB:         spdxcommon.DocElementID{ElementRefID: pkgID},
			Relationship: "CONTAINS",
		})
	}

	// TODO: https://github.com/distix-pj/distix/issues/14
	for _, pkg := range sys.Packages {
		pkgID := spdxcommon.ElementID(fmt.Sprintf("SPDXRef-Pkg-%s", pkg.GetPurl()))
		for _, req := range pkg.Requires {
			if reqID, ok := providerMap[req.Name]; ok {
				relationships = append(relationships, &spdxv2_3.Relationship{
					RefA:         spdxcommon.DocElementID{ElementRefID: pkgID},
					RefB:         spdxcommon.DocElementID{ElementRefID: reqID},
					Relationship: "DEPENDS_ON",
				})
			}
			// else {
			// }
		}
	}

	doc := &spdxv2_3.Document{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXIdentifier:    "DOCUMENT",
		DocumentName:      sys.HostName,
		DocumentNamespace: "https://distix.example.org/package/" + sys.HostName,
		CreationInfo: &spdxv2_3.CreationInfo{
			Created: time.Now().UTC().Format(time.RFC3339),
			Creators: []spdxcommon.Creator{
				{CreatorType: "Tool", Creator: "distix"},
			},
		},
		Packages:          packages,
		Relationships:     relationships,
		// Files:             files,
	}

	return w.write(doc, out)
}

func (w *SPDXWriter) WriteDistSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
}


func (w *SPDXWriter) write(doc *spdxv2_3.Document, out io.Writer) error {
	switch w.fileFormat {
	case JSON:
		return spdxjson.Write(doc, out, spdxjson.Indent("\t"))
	default:
		return fmt.Errorf("SPDX file format %s is not yet supported", w.fileFormat)
	}
}


func toSPDXAlgorithm(algo model.DigestAlgorithm) (spdxcommon.ChecksumAlgorithm, error) {
	switch algo {
	case model.DigestAlgorithmMD5:
		return spdxcommon.MD5, nil
	case model.DigestAlgorithmSHA1:
		return spdxcommon.SHA1, nil
	case model.DigestAlgorithmSHA256:
		return spdxcommon.SHA256, nil
	case model.DigestAlgorithmSHA384:
		return spdxcommon.SHA384, nil
	case model.DigestAlgorithmSHA512:
		return spdxcommon.SHA512, nil
	case model.DigestAlgorithmSHA224:
		return spdxcommon.SHA224, nil
	case model.DigestAlgorithmMD2:
		return spdxcommon.MD2, nil
	default:
		// TODO: implement Stringer for human-readable algorithm names in error messages
		return "", fmt.Errorf("unsupported digest algorithm: %d", algo)
	}
}

