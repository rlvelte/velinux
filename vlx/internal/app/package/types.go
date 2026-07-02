package _package

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// TODO: Maybe use more info about packages later...

type PackageType string

const (
	PackageTypePackage PackageType = "package"
	PackageTypePattern PackageType = "pattern"
	PackageTypePatch   PackageType = "patch"
	PackageTypeProduct PackageType = "product"
)

type Package struct {
	Name        string
	Version     string
	Arch        string
	Repo        string
	Installed   bool
	Upgradable  bool
	Type        PackageType
	Description string
}

type Stream struct {
	XMLName      xml.Name      `xml:"stream"`
	SearchResult *SearchResult `xml:"search-result"`
}

type SearchResult struct {
	Version      string       `xml:"version,attr"`
	SolvableList SolvableList `xml:"solvable-list"`
}

type SolvableList struct {
	XMLName   xml.Name   `xml:"solvable-list"`
	Solvables []Solvable `xml:"solvable"`
}

type Solvable struct {
	Status      string `xml:"status,attr"`
	Name        string `xml:"name,attr"`
	Kind        string `xml:"kind,attr"`
	Edition     string `xml:"edition,attr"`
	Arch        string `xml:"arch,attr"`
	Repository  string `xml:"repository,attr"`
	Description string `xml:"description"`
}

func parseXml(data []byte) ([]Package, error) {
	var stream Stream
	if err := xml.Unmarshal(data, &stream); err != nil {
		return nil, fmt.Errorf("failed to parse zypper search XML: %w", err)
	}

	if stream.SearchResult == nil {
		return []Package{}, nil
	}

	solvables := stream.SearchResult.SolvableList.Solvables
	pkgs := make([]Package, 0, len(solvables))

	for _, s := range solvables {
		edition := s.Edition
		var version, release string

		if idx := strings.Index(edition, "-"); idx >= 0 {
			version = edition[:idx]
			release = edition[idx+1:]
		} else {
			version = edition
		}

		_ = release

		kind := s.Kind
		switch strings.ToLower(kind) {
		case "paket", "package":
			kind = "package"
		case "muster", "pattern":
			kind = "pattern"
		case "patch":
			kind = "patch"
		case "produkt", "product":
			kind = "product"
		}

		pkg := Package{
			Name:        s.Name,
			Version:     version,
			Arch:        s.Arch,
			Repo:        s.Repository,
			Installed:   s.Status == "installed",
			Upgradable:  s.Status == "upgradable",
			Type:        PackageType(kind),
			Description: s.Description,
		}

		pkgs = append(pkgs, pkg)
	}

	return pkgs, nil
}
