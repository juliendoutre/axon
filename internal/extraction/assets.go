package extraction

import "regexp"

type Asset struct {
	Type string
	ID   string
}

var SupportedAssets = map[string]regexp.Regexp{} //nolint:gochecknoglobals

func ListMatches(candidate Candidate) []Asset {
	assets := []Asset{}

	for assetType, regex := range SupportedAssets {
		for _, match := range regex.FindAllString(candidate.Value, -1) {
			assets = append(assets, Asset{Type: assetType, ID: match})
		}
	}

	return assets
}
