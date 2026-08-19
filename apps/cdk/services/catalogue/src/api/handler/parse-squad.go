package main

import (
	"catalogue/src/types/dynamo"
	"contracts/dist/catalogue/v1"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func parseDynamoSquad(item *dynamo.Squad) *catalogue.Squad {
	loadouts := []*catalogue.Loadout{}
	for _, loadout := range item.Loadouts {
		if unit, found := Units[loadout.Unit]; found {
			var item *catalogue.Item
			if loadout.Item != nil {
				if it, found := Items[*loadout.Item]; found {
					item = it
				}
			}
			loadouts = append(loadouts, &catalogue.Loadout{
				Unit: unit,
				Item: item,
			})
		}
	}

	return &catalogue.Squad{
		Id:           item.Id,
		Version:      item.Version,
		CreatedTime:  timestamppb.New(time.UnixMilli(item.CreatedTime)),
		ModifiedTime: timestamppb.New(time.UnixMilli(item.ModifiedTime)),

		Name:     item.Name,
		Faction:  item.Faction,
		Loadouts: loadouts,
	}
}
