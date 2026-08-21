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
	"NuCZhM3zSjLZuqCtYwCX": {
		Id:           "NuCZhM3zSjLZuqCtYwCX",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:        new("/clone-scout.png"),
		Name:         "Clone Scout",
		Price:        180,
		Size:         1,
		Speed:        2,
		Health:       15,
		ArmorClass:   3,
		Marksmanship: new(uint32(2)),
		MeleeAbility: new(uint32(2)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["Kq0BO0W0rH2Ko0kyydLI"], // Blaster
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			{
				Id:      "bx2mnRh3jbPrkWUTtExU",
				Version: 1,
				Name:    "Unarmed Strike",
				Details: &catalogue.Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   6,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    8,
							Modifier: -2,
						},
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["AoM09AjZL8hUMhTYLIYY"], // Vibroblade
		},
	},
	"vEPHxbpBSxzeNluR4b9U": {
		Id:           "vEPHxbpBSxzeNluR4b9U",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:        new("/clone-sharpshooter.png"),
		Name:         "Clone Sharpshooter",
		Price:        230,
		Size:         1,
		Speed:        2,
		Health:       15,
		ArmorClass:   3,
		Marksmanship: new(uint32(3)),
		MeleeAbility: new(uint32(2)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["UG60YI7qn8DJryG9pmR7"], // Blaster Rifle
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			{
				Id:      "MS0sOPe1hYsbegxmLI5V",
				Version: 1,
				Name:    "Unarmed Strike",
				Details: &catalogue.Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   6,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    8,
							Modifier: -1,
						},
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["AoM09AjZL8hUMhTYLIYY"], // Vibroblade
		},
	},
	"AFFXxbwBD1aNFYbt25m7": {
		Id:           "AFFXxbwBD1aNFYbt25m7",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:        new("/clone-commander.png"),
		Name:         "Clone Commander",
		Price:        250,
		Size:         1,
		Speed:        2,
		Health:       15,
		ArmorClass:   3,
		Marksmanship: new(uint32(3)),
		MeleeAbility: new(uint32(3)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["OiFq73Vk00ob7p8IuKAb"], // Hand Blasters
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			{
				Id:      "KZYHYYxMJznfWqDn7lTB",
				Version: 1,
				Name:    "Unarmed Strike",
				Details: &catalogue.Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   6,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    8,
							Modifier: -1,
						},
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["AoM09AjZL8hUMhTYLIYY"], // Vibroblade
		},
	},
	"buMbZfN6zRflygTOVHaG": {
		Id:           "buMbZfN6zRflygTOVHaG",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_GALACTIC_REPUBLIC},
		Image:        new("/fighter-tank.png"),
		Name:         "Fighter Tank",
		Price:        480,
		Size:         4,
		Speed:        4,
		Health:       30,
		ArmorClass:   5,
		Marksmanship: new(uint32(3)),
		Actions:      []*catalogue.Unit_Action{},
		Items: []*catalogue.Item{
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
			{
				Id:      "VS8ICW50OVr0L5iPxS0Y",
				Version: 1,
				Name:    "Laser Cannons",
				Details: &catalogue.Item_Cannon{
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
				Id:      "FC9rHpwti4FIuyhH6bs1",
				Version: 1,
				Name:    "Rocket Launchers",
				Details: &catalogue.Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 2,
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
		},
	},
	"5LDT5irCFiLuKt6wQMvo": {
		Id:           "5LDT5irCFiLuKt6wQMvo",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:        new("/droid-scout.png"),
		Name:         "Droid",
		Price:        80,
		Size:         1,
		Speed:        2,
		Health:       10,
		ArmorClass:   1,
		Marksmanship: new(uint32(1)),
		MeleeAbility: new(uint32(1)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["Kq0BO0W0rH2Ko0kyydLI"], // Blaster
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			{
				Id:      "6ZMSvxmGw7zO1sWeJbig",
				Version: 1,
				Name:    "Unarmed Strike",
				Details: &catalogue.Item_MeleeWeapon{
					MeleeWeapon: &catalogue.MeleeWeapon{
						AttackSpeed:   5,
						ArmorPiercing: 1,
						Damage: &common.DiceRoll{
							Count:    1,
							Sides:    8,
							Modifier: -3,
						},
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"6EpGE8Td4gLWkcYgDIW9": {
		Id:           "6EpGE8Td4gLWkcYgDIW9",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:        new("/super-battle-droid.png"),
		Name:         "Super Battle Droid",
		Price:        170,
		Size:         1,
		Speed:        1,
		Health:       17,
		ArmorClass:   2,
		Marksmanship: new(uint32(2)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			{
				Id:      "Vgaa0TwEqbfHYRnGrkdf",
				Version: 1,
				Name:    "Arm blasters",
				Details: &catalogue.Item_Blaster{
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
		LoadoutItems: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"R4ov67MKvT4YFauIDJnV": {
		Id:           "R4ov67MKvT4YFauIDJnV",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:        new("/dwarf-spider-droid.png"),
		Name:         "Dwarf Spider Droid",
		Price:        370,
		Size:         2,
		Speed:        3,
		Health:       20,
		ArmorClass:   2,
		Marksmanship: new(uint32(3)),
		Actions: []*catalogue.Unit_Action{
			{
				Name:        "Sprint",
				Description: new("Double move"),
			},
		},
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			{
				Id:      "Z5WYIKT9RLCj6VQ6TfZe",
				Version: 1,
				Name:    "Blaster Rifle",
				Details: &catalogue.Item_Blaster{
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
		LoadoutItems: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"19eF4NtFTiuBXPCLFCNU": {
		Id:           "19eF4NtFTiuBXPCLFCNU",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:        new("/octuptarra.png"),
		Name:         "Octuptarra",
		Price:        460,
		Size:         4,
		Speed:        3,
		Health:       25,
		ArmorClass:   3,
		Marksmanship: new(uint32(3)),
		Actions:      []*catalogue.Unit_Action{},
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
			{
				Id:      "mvOUbhmAkRnJtRVDrRrs",
				Version: 1,
				Name:    "Laser Blaster",
				Details: &catalogue.Item_Blaster{
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
				Id:      "thVjEWWJSKBxtcIo45uE",
				Version: 1,
				Name:    "Rocket Launcher",
				Details: &catalogue.Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 1,
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
		},
	},
	"7JSXGhVxkJUuNjZFp3KY": {
		Id:           "7JSXGhVxkJUuNjZFp3KY",
		Version:      1,
		Factions:     []catalogue.Faction{catalogue.Faction_FACTION_SEPARATIST_ALLIANCE},
		Image:        new("/armored-assault-tank.png"),
		Name:         "Armored Assault Tank",
		Price:        530,
		Size:         4,
		Speed:        3,
		Health:       30,
		ArmorClass:   4,
		Marksmanship: new(uint32(2)),
		Actions:      []*catalogue.Unit_Action{},
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
			{
				Id:      "j2QtTUcX1xGojMhOTExc",
				Version: 1,
				Name:    "Side Blasters",
				Details: &catalogue.Item_Blaster{
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
				Id:      "y7hOdf9HS1BqNz59yCul",
				Version: 1,
				Name:    "Main Turret",
				Details: &catalogue.Item_Cannon{
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
				Id:      "o6sx09eU7AWx6F0b0T8e",
				Version: 1,
				Name:    "Rocket Launchers",
				Details: &catalogue.Item_Launcher{
					Launcher: &catalogue.Launcher{
						FireRate: 2,
					},
				},
			},
		},
		LoadoutItems: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
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
