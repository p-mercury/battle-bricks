package dynamo

import (
	"contracts/dist/catalogue/v1"
)

type Squad struct {
	Pk     string
	Sk     string
	Gsi1pk string
	Gsi1sk string
	Type   string

	Id           string
	Correlations map[string]Correlation
	Version      uint64
	CreatedTime  int64
	ModifiedTime int64

	Name     string
	Faction  catalogue.Faction
	Loadouts []struct {
		Unit string
		Item *string
	}
}
