package format

import (
	"fmt"
	"io"

	"github.com/distix-pj/distix/data/model"
	"github.com/distix-pj/distix/format/cyclonedx"
	"github.com/distix-pj/distix/format/spdx"
	"github.com/distix-pj/distix/format/types"
)

type Writer interface {
	WritePackage(*model.Package, io.Writer) error
	WriteOneSystem(*model.System, io.Writer) error
	WriteDistSystem(*model.System, io.Writer) error
}

func NewWriter(sbomType types.SbomType) (Writer, error) {
	switch sbomType.RecordType {
	case types.SPDX:
		return spdx.NewWriter(sbomType.FileFormatType),nil
	case types.CYCLONEDX:
		return cyclonedx.NewWriter(sbomType.FileFormatType),nil
	default:
		return nil, fmt.Errorf("unsupported SBOM record type: %s", sbomType.RecordType)
	}
}

