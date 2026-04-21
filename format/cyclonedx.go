package format

import (
	"errors"
	"fmt"
	"io"

	cdx "github.com/CycloneDX/cyclonedx-go"

	"github.com/distix-pj/distix/data/model"
)

type CDXWriter struct {
	fileFormat SbomFileFormatType
}

func (w *CDXWriter) WritePackage(pkg *model.Package, out io.Writer) error {
	mainComp := cdx.Component{
		Type:       cdx.ComponentTypeLibrary,
		Name:       pkg.PkgNevra.Name,
		Version:    pkg.PkgNevra.Version,
		PackageURL: pkg.GetPurl(),
		BOMRef:     pkg.GetPurl(),
	}

	components := []cdx.Component{}
	depRefs := []string{}

	for i, prov := range pkg.Provides {
		provRef := fmt.Sprintf("RPM-PROV-%d", i)
		components = append(components, cdx.Component{
			Type:       cdx.ComponentTypeLibrary,
			Name:       prov.Name,
			Version:    prov.Version,
			BOMRef:     provRef,
		})
	}

	for i, req := range pkg.Requires {
		reqRef := fmt.Sprintf("RPM-CAP-%d", i)
		components = append(components, cdx.Component{
			Type:       cdx.ComponentTypeLibrary,
			Name:       req.Name,
			Version:    req.Version,
			BOMRef:     reqRef,
		})
		depRefs = append(depRefs, reqRef)
	}

	for i, file := range pkg.Files {
		fileRef := fmt.Sprintf("RPM-FILE-%d", i)
		algo, err := toCDXAlgorithm(file.DigestAlgorithm)
		if err != nil {
			return err
		}
		components = append(components, cdx.Component{
			Type:       cdx.ComponentTypeFile,
			Name:       file.Name,
			BOMRef:     fileRef,
			Hashes: &[]cdx.Hash{{
				Algorithm: algo,
				Value:     file.Digest,
			}},
		})
	}

	deps := []cdx.Dependency{
		{Ref: pkg.GetPurl(), Dependencies: &depRefs},
	}

	bom := cdx.NewBOM()
	bom.Metadata = &cdx.Metadata{Component: &mainComp}
	bom.Components = &components
	bom.Dependencies = &deps

	return w.encode(bom, out)
}

func (w *CDXWriter) WriteOneSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
}

func (w *CDXWriter) WriteDistSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
}

func (w *CDXWriter) encode(bom *cdx.BOM, out io.Writer) error {
	var format cdx.BOMFileFormat
	switch w.fileFormat {
	case JSON:
		format = cdx.BOMFileFormatJSON
	default:
		return fmt.Errorf("CycloneDX file format %s is not supported", w.fileFormat)
	}
	enc := cdx.NewBOMEncoder(out, format)
	enc.SetPretty(true)
	return enc.Encode(bom)
}


func toCDXAlgorithm(algo model.DigestAlgorithm) (cdx.HashAlgorithm, error) {
	switch algo {
	case model.DigestAlgorithmMD5:
		return cdx.HashAlgoMD5, nil
	case model.DigestAlgorithmSHA1:
		return cdx.HashAlgoSHA1, nil
	case model.DigestAlgorithmSHA256:
		return cdx.HashAlgoSHA256, nil
	case model.DigestAlgorithmSHA384:
		return cdx.HashAlgoSHA384, nil
	case model.DigestAlgorithmSHA512:
		return cdx.HashAlgoSHA512, nil
	default:
		// TODO: implement Stringer for human-readable algorithm names in error messages
		return "", fmt.Errorf("unsupported digest algorithm: %d", algo)
	}
}

