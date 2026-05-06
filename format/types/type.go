package types

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type SbomRecordType string
const (
    SPDX      SbomRecordType = "spdx"
    CYCLONEDX SbomRecordType = "cyclonedx"
)

func SbomRecordTypeFromStr(s string) (SbomRecordType, error) {
	switch SbomRecordType(s) {
	case SPDX, CYCLONEDX:
		return SbomRecordType(s), nil
	default:
		return "", fmt.Errorf("unknown SbomRecordType: %s", s)
	}
}

func (s *SbomRecordType) String() string {
	return string(*s)
}
func (s *SbomRecordType) Type() string {
	return "SbomRecordType"
}



type SbomFileFormatType string
const (
    JSON     SbomFileFormatType = "json"
    TAGVALUE SbomFileFormatType = "tagvalue"
    XML      SbomFileFormatType = "xml"
    YAML     SbomFileFormatType = "yaml"
)

func SbomFileFormatTypeFromStr(s string) (SbomFileFormatType, error) {
	switch SbomFileFormatType(s) {
	case JSON, TAGVALUE, XML, YAML:
		return SbomFileFormatType(s), nil
	default:
		return "", fmt.Errorf("unknown SbomFileFormatType: %s", s)
	}
}

func (s *SbomFileFormatType) String() string {
    return string(*s)
}
func (s *SbomFileFormatType) Type() string {
	return "SbomFileFormatType"
}



type SbomType struct {
	RecordType     SbomRecordType
	FileFormatType SbomFileFormatType
}

// TODO: https://github.com/distix-pj/distix/issues/15
var ValidSbomTypes = map[SbomRecordType][]SbomFileFormatType{
	SPDX: {
		JSON,
		XML,
		YAML,
		TAGVALUE,
	},
	CYCLONEDX: {
		JSON,
		XML,
	},
}

func NewSbomType(recordType SbomRecordType, fileFormatType SbomFileFormatType) (*SbomType, error) {
	sbomType := &SbomType{
		RecordType:     recordType,
		FileFormatType: fileFormatType,
	}
	if !sbomType.IsValid() {
		return nil, errors.New("Invalid SbomType")
	}
	return sbomType, nil
}

func (s *SbomType) IsValid() bool {
	formatTypes, exists := ValidSbomTypes[s.RecordType]
	if !exists {
		return false
	}

	for _, format := range formatTypes {
		if format == s.FileFormatType {
			return true
		}
	}
	return false
}

func (s *SbomType) String() string {
	if s.RecordType == "" || s.FileFormatType == "" {
		return ""
	}
	return string(s.RecordType) + "-" + string(s.FileFormatType)
}

func (s *SbomType) Set(str string) error {
	index := strings.Index(str, "-")
	if index <= 0 {
		return errors.New("Invalid SbomType")
	}
	sbomRecordType, err := SbomRecordTypeFromStr(str[:index])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	sbomFileFormatType, err := SbomFileFormatTypeFromStr(str[index+1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	s.RecordType = sbomRecordType
	s.FileFormatType = sbomFileFormatType
	return nil
}

func (s *SbomType) Type() string {
	return "SbomType"
}


func GetSbomTypeDefault() SbomType {
	defSbomType, _ := NewSbomType(SPDX, JSON)
	return *defSbomType
}
