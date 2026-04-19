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
	}

	return w.write(doc, out)
}

func (w *SPDXWriter) WriteOneSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
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

