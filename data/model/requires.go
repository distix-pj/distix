package model

import (
	"fmt"
	"regexp"
)


type RpmCapability struct {
	Name    string
	Version string
	Flags   uint32
}

func (rCap *RpmCapability) IsSoLib() bool {
	rpmSoPattern := regexp.MustCompile(`^.*\.so(\.\d+)*(\([^)]*\))+$`)
	return rpmSoPattern.MatchString(rCap.Name)
}

func (rCap *RpmCapability) GetId() string {
	return fmt.Sprintf("RPM-CAP-%s", rCap.Name)
}


// DigestAlgorithm represents hash algorithm identifiers as defined in RFC 4880.
type DigestAlgorithm int
const (
	DigestAlgorithmUnknown  DigestAlgorithm = 0
	DigestAlgorithmMD5      DigestAlgorithm = 1
	DigestAlgorithmSHA1     DigestAlgorithm = 2
	DigestAlgorithmRIPEMD160 DigestAlgorithm = 3
	DigestAlgorithmMD2      DigestAlgorithm = 5
	DigestAlgorithmTIGER192 DigestAlgorithm = 6
	DigestAlgorithmHAVAL5160 DigestAlgorithm = 7
	DigestAlgorithmSHA256   DigestAlgorithm = 8
	DigestAlgorithmSHA384   DigestAlgorithm = 9
	DigestAlgorithmSHA512   DigestAlgorithm = 10
	DigestAlgorithmSHA224   DigestAlgorithm = 11
)


type RpmFile struct {
	Name string
	Digest string
	DigestAlgorithm DigestAlgorithm
	// fileInfo rpmtuils.FileInfo
}

func (rFile *RpmFile) GetId() string {
	return fmt.Sprintf("RPM-FILE-%s", rFile.Name)
}

