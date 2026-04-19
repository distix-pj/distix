package format

import (
	"errors"
	"io"

	"github.com/distix-pj/distix/data/model"
)

type CDXWriter struct {
	fileFormat SbomFileFormatType
}

func (w *CDXWriter) WritePackage(pkg *model.Package, out io.Writer) error {
	return errors.New("not implemented")
}

func (w *CDXWriter) WriteOneSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
}

func (w *CDXWriter) WriteDistSystem(sys *model.System, out io.Writer) error {
	return errors.New("not implemented")
}

