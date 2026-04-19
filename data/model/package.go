package model

import (
	"fmt"
)


type Package struct {
	PkgNevra PackageNevra
	Summary	string
	Description	string
	Provides []RpmCapability
	Requires []RpmCapability
	Files []RpmFile

	//! used by onesystem
	RequirePkgs []string
}


func (pkg *Package) getNEVRA() string {
	return pkg.PkgNevra.GetNEVRA()
}
func (pkg *Package) getName() string {
	return pkg.PkgNevra.Name
}
func (pkg *Package) getVersion() string {
	return pkg.PkgNevra.Version
}
func (pkg *Package) getVersionRelease() string {
	return pkg.PkgNevra.GetVR()
}
//! TODO: This is instant implementation. Need to be fixed.
func (pkg *Package) GetPurl() string {
	return fmt.Sprintf("pkg:rpm/generic/%s@%s", pkg.getName(), pkg.getVersionRelease())
}

