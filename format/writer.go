package format

import (
	"fmt"
	"io"

	"github.com/distix-pj/distix/data/model"
)

type Writer interface {
	WritePackage(*model.Package, io.Writer) error
	WriteOneSystem(*model.System, io.Writer) error
	WriteDistSystem(*model.System, io.Writer) error
}

func NewWriter(sbomType SbomType) (Writer, error) {
	switch sbomType.RecordType {
	case SPDX:
		return &SPDXWriter{fileFormat: sbomType.FileFormatType},nil
	case CYCLONEDX:
		return &CDXWriter{fileFormat: sbomType.FileFormatType},nil
	default:
		return nil, fmt.Errorf("unsupported SBOM record type: %s", sbomType.RecordType)
	}
}

