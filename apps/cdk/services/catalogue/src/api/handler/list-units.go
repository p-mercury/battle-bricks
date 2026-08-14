package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"contracts/dist/common/v1"
	"log/slog"

	"connectrpc.com/connect"
)

var Units = map[string]*catalogue.Unit{
	"vEPHxbpBSxzeNluR4b9U": {
		Id:            "vEPHxbpBSxzeNluR4b9U",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:         new("/unit-clone.png"),
		Name:          "Clone Trooper",
		Price:         150,
		Size:          1,
		Speed:         2,
		Health:        14,
		ArmorClass:    3,
		CarryCapacity: 15,
		Marksmanship:  new(uint32(2)),
		MeleeAbility:  new(uint32(2)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Unarmed Strike",
				Details: &catalogue.Unit_Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   4,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    6,
							Modifier: -3,
						},
					},
				},
			},
		},
	},
	"AFFXxbwBD1aNFYbt25m7": {
		Id:            "AFFXxbwBD1aNFYbt25m7",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:         new("/unit-clone-specialist.png"),
		Name:          "Clone Specialist",
		Price:         190,
		Size:          1,
		Speed:         2,
		Health:        14,
		ArmorClass:    3,
		CarryCapacity: 15,
		Marksmanship:  new(uint32(3)),
		MeleeAbility:  new(uint32(3)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Unarmed Strike",
				Details: &catalogue.Unit_Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   4,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    6,
							Modifier: -1,
						},
					},
				},
			},
		},
	},
	"buMbZfN6zRflygTOVHaG": {
		Id:            "buMbZfN6zRflygTOVHaG",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:         new("/fighter-tank.png"),
		Name:          "Fighter Tank",
		Price:         350,
		Size:          4,
		Speed:         4,
		Health:        30,
		ArmorClass:    4,
		CarryCapacity: 24,
		Marksmanship:  new(uint32(3)),
		Actions:       []*catalogue.Unit_Action{},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Laser Cannons",
				Details: &catalogue.Unit_Item_Cannon{
					Cannon: &catalogue.Cannon{
						Range: &catalogue.Range{
							Min: 5,
							Max: 12,
						},
						FireRate: 2,
					},
				},
			},
			{
				Name: "Rocket Launchers",
				Details: &catalogue.Unit_Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 2,
					},
				},
			},
		},
	},
	"5LDT5irCFiLuKt6wQMvo": {
		Id:            "5LDT5irCFiLuKt6wQMvo",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:         new("/unit-droid.png"),
		Name:          "Droid",
		Price:         55,
		Size:          1,
		Speed:         2,
		Health:        10,
		ArmorClass:    1,
		CarryCapacity: 15,
		Marksmanship:  new(uint32(1)),
		MeleeAbility:  new(uint32(1)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Unarmed Strike",
				Details: &catalogue.Unit_Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   3,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    6,
							Modifier: -3,
						},
					},
				},
			},
		},
	},
	"6EpGE8Td4gLWkcYgDIW9": {
		Id:            "6EpGE8Td4gLWkcYgDIW9",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:         new("/unit-super-battle-droid.png"),
		Name:          "Super Battle Droid",
		Price:         160,
		Size:          1,
		Speed:         1,
		Health:        14,
		ArmorClass:    2,
		CarryCapacity: 9,
		Marksmanship:  new(uint32(2)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Arm blasters",
				Details: &catalogue.Unit_Item_Blaster{
					Blaster: &catalogue.Blaster{
						Range: &catalogue.Range{
							Min: 0,
							Max: 8,
						},
						FireRate: 10,
					},
				},
			},
		},
	},
	"R4ov67MKvT4YFauIDJnV": {
		Id:            "R4ov67MKvT4YFauIDJnV",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:         new("/dwarf-spider-droid.png"),
		Name:          "Dwarf Spider Droid",
		Price:         360,
		Size:          2,
		Speed:         3,
		Health:        20,
		ArmorClass:    2,
		CarryCapacity: 24,
		Marksmanship:  new(uint32(3)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Blaster Rifle",
				Details: &catalogue.Unit_Item_Blaster{
					Blaster: &catalogue.Blaster{
						Range: &catalogue.Range{
							Min: 2,
							Max: 14,
						},
						FireRate: 8,
					},
				},
			},
		},
	},
	"19eF4NtFTiuBXPCLFCNU": {
		Id:            "19eF4NtFTiuBXPCLFCNU",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:         new("/octuptarra.png"),
		Name:          "Octuptarra",
		Price:         360,
		Size:          4,
		Speed:         3,
		Health:        22,
		ArmorClass:    3,
		CarryCapacity: 24,
		Marksmanship:  new(uint32(3)),
		Actions:       []*catalogue.Unit_Action{},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Laser Blaster",
				Details: &catalogue.Unit_Item_Blaster{
					Blaster: &catalogue.Blaster{
						Range: &catalogue.Range{
							Min: 4,
							Max: 10,
						},
						FireRate: 10,
					},
				},
			},
			{
				Name: "Rocket Launcher",
				Details: &catalogue.Unit_Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 1,
					},
				},
			},
		},
	},
	"7JSXGhVxkJUuNjZFp3KY": {
		Id:            "7JSXGhVxkJUuNjZFp3KY",
		Version:       1,
		Factions:      []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:         new("/armored-assault-tank.png"),
		Name:          "Armored Assault Tank",
		Price:         400,
		Size:          4,
		Speed:         3,
		Health:        26,
		ArmorClass:    4,
		CarryCapacity: 24,
		Marksmanship:  new(uint32(2)),
		Actions:       []*catalogue.Unit_Action{},
		Items: []*catalogue.Unit_Item{
			{
				Name: "Side Blasters",
				Details: &catalogue.Unit_Item_Blaster{
					Blaster: &catalogue.Blaster{
						Range: &catalogue.Range{
							Min: 4,
							Max: 10,
						},
						FireRate: 10,
					},
				},
			},
			{
				Name: "Main Turret",
				Details: &catalogue.Unit_Item_Cannon{
					Cannon: &catalogue.Cannon{
						Range: &catalogue.Range{
							Min: 5,
							Max: 12,
						},
						FireRate: 1,
					},
				},
			},
			{
				Name: "Rocket Launchers",
				Details: &catalogue.Unit_Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 2,
					},
				},
			},
		},
	},
}

func (s *Handler) ListUnits(
	ctx context.Context,
	req *connect.Request[catalogue.ListUnitsRequest],
) (*connect.Response[catalogue.ListUnitsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_units",
			&catalogue.ListUnitsContext{
				Request: req.Msg,
			},
		)
		if err != nil {
			logger.Error("Error evaluating policy", slog.Any("error", err))
			return nil, connectkit.NewUnexpected()
		} else if !resp.Authz {
			return nil, connectkit.NewUnauthorized()
		}
	}

	values := make([]*catalogue.Unit, 0, len(Units))
	for _, value := range Units {
		values = append(values, value)
	}

	return connect.NewResponse(&catalogue.ListUnitsResponse{Units: values}), nil
}
