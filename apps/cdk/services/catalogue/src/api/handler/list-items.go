package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"contracts/dist/common/v1"
	"log/slog"

	"connectrpc.com/connect"
)

var Items = map[string]*catalogue.Item{
	"CV1AILTJNoyVKD7Uon63": {
		Id:          "CV1AILTJNoyVKD7Uon63",
		Version:     1,
		Name:        "Red Plasma Cartridge",
		Description: new("Cheap Tibanna gas blend producing weaker red bolts"),
		Price:       5,
		Weight:      2,
		Details: &catalogue.Item_BlasterBolt{
			BlasterBolt: &catalogue.BlasterBolt{
				Capacity:      40,
				ArmorPiercing: 1,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: -1,
				},
			},
		},
	},
	"d2D9WroLYwyxAWaUpG1u": {
		Id:          "d2D9WroLYwyxAWaUpG1u",
		Version:     1,
		Name:        "Blue Plasma Cartridge",
		Description: new("High grade Tibanna gas blend producing blue bolts"),
		Price:       10,
		Weight:      2,
		Details: &catalogue.Item_BlasterBolt{
			BlasterBolt: &catalogue.BlasterBolt{
				Capacity:      35,
				ArmorPiercing: 2,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 0,
				},
			},
		},
	},
	"qL0jBLWBgS6D5mXTcxF5": {
		Id:          "qL0jBLWBgS6D5mXTcxF5",
		Version:     1,
		Name:        "Green Plasma Cartridge",
		Description: new("Pure refined Tibanna producing powerful green bolts"),
		Price:       15,
		Weight:      2,
		Details: &catalogue.Item_BlasterBolt{
			BlasterBolt: &catalogue.BlasterBolt{
				Capacity:      30,
				ArmorPiercing: 2,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 1,
				},
			},
		},
	},
	"QMdY2xQJD9UEcLDwTNTp": {
		Id:          "QMdY2xQJD9UEcLDwTNTp",
		Version:     1,
		Name:        "Yellow Plasma Cartridge",
		Description: new("High pressure Tibanna producing armor piercing yellow bolts"),
		Price:       20,
		Weight:      2,
		Details: &catalogue.Item_BlasterBolt{
			BlasterBolt: &catalogue.BlasterBolt{
				Capacity:      25,
				ArmorPiercing: 3,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 1,
				},
			},
		},
	},
	"Agkg3j4cpGdmzIda67jD": {
		Id:          "Agkg3j4cpGdmzIda67jD",
		Version:     1,
		Name:        "Red Plasma Shell",
		Description: new("Cheap Tibanna gas blend producing weaker blasts"),
		Price:       10,
		Weight:      4,
		Details: &catalogue.Item_CannonShell{
			CannonShell: &catalogue.CannonShell{
				SplashRadius:  0,
				ArmorPiercing: 3,
				Damage: &common.DiceRoll{
					Count:    2,
					Sides:    8,
					Modifier: 0,
				},
			},
		},
	},
	"bJeepLfSKSzI1yznvTUM": {
		Id:          "bJeepLfSKSzI1yznvTUM",
		Version:     1,
		Name:        "Blue Plasma Shell",
		Description: new("High grade Tibanna gas blend producing blue bolts"),
		Price:       20,
		Weight:      4,
		Details: &catalogue.Item_CannonShell{
			CannonShell: &catalogue.CannonShell{
				SplashRadius:  0,
				ArmorPiercing: 3,
				Damage: &common.DiceRoll{
					Count:    2,
					Sides:    8,
					Modifier: 2,
				},
			},
		},
	},
	"bPC0uxvBFkmHWaHDrHoB": {
		Id:          "bPC0uxvBFkmHWaHDrHoB",
		Version:     1,
		Name:        "Green Plasma Shell",
		Description: new("Pure refined Tibanna producing powerful green bolts"),
		Price:       30,
		Weight:      4,
		Details: &catalogue.Item_CannonShell{
			CannonShell: &catalogue.CannonShell{
				SplashRadius:  0,
				ArmorPiercing: 4,
				Damage: &common.DiceRoll{
					Count:    2,
					Sides:    8,
					Modifier: 4,
				},
			},
		},
	},
	"xl4SlN7caFSq68DPxhr3": {
		Id:          "xl4SlN7caFSq68DPxhr3",
		Version:     1,
		Name:        "Fragmentation rocket",
		Description: new("Anti-personnel rocket that shreds lightly armored targets with a wide blast radius"),
		Price:       30,
		Weight:      6,
		Details: &catalogue.Item_LauncherRocket{
			LauncherRocket: &catalogue.LauncherRocket{
				Range: &catalogue.Range{
					Min: 6,
					Max: 14,
				},
				SplashRadius:  2,
				ArmorPiercing: 2,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 4,
				},
			},
		},
	},
	"2A06zTJVEAHqG8Ax4w7k": {
		Id:          "2A06zTJVEAHqG8Ax4w7k",
		Version:     1,
		Name:        "Ion rocket",
		Description: new("Anti-vehicle rocket for piercing thick armor"),
		Price:       30,
		Weight:      6,
		Details: &catalogue.Item_LauncherRocket{
			LauncherRocket: &catalogue.LauncherRocket{
				Range: &catalogue.Range{
					Min: 6,
					Max: 14,
				},
				SplashRadius:  0,
				ArmorPiercing: 4,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 4,
				},
			},
		},
	},
	"Kq0BO0W0rH2Ko0kyydLI": {
		Id:      "Kq0BO0W0rH2Ko0kyydLI",
		Version: 1,
		Name:    "Blaster",
		Price:   20,
		Weight:  4,
		Details: &catalogue.Item_Blaster{
			Blaster: &catalogue.Blaster{
				Range: &catalogue.Range{
					Min: 2,
					Max: 12,
				},
				FireRate: 8,
			},
		},
	},
	"OiFq73Vk00ob7p8IuKAb": {
		Id:      "OiFq73Vk00ob7p8IuKAb",
		Version: 1,
		Name:    "Hand Blasters",
		Price:   35,
		Weight:  5,
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
	"UG60YI7qn8DJryG9pmR7": {
		Id:      "UG60YI7qn8DJryG9pmR7",
		Version: 1,
		Name:    "Blaster Rifle",
		Price:   30,
		Weight:  5,
		Details: &catalogue.Item_Blaster{
			Blaster: &catalogue.Blaster{
				Range: &catalogue.Range{
					Min: 5,
					Max: 20,
				},
				FireRate: 6,
			},
		},
	},
	"AoM09AjZL8hUMhTYLIYY": {
		Id:      "AoM09AjZL8hUMhTYLIYY",
		Version: 1,
		Name:    "Vibroblade",
		Price:   15,
		Weight:  1,
		Details: &catalogue.Item_MeleeWeapon{
			MeleeWeapon: &catalogue.MeleeWeapon{
				AttackSpeed:   5,
				ArmorPiercing: 2,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 0,
				},
			},
		},
	},
	"sxBkc0uTTc66u2Bi8fG1": {
		Id:          "sxBkc0uTTc66u2Bi8fG1",
		Version:     1,
		Name:        "",
		Description: new(""),
		Price:       25,
		Weight:      2,
		Details: &catalogue.Item_MeleeWeapon{
			MeleeWeapon: &catalogue.MeleeWeapon{
				AttackSpeed:   4,
				ArmorPiercing: 3,
				Damage: &common.DiceRoll{
					Count:    1,
					Sides:    8,
					Modifier: 2,
				},
			},
		},
	},
}

func (s *Handler) ListItems(
	ctx context.Context,
	req *connect.Request[catalogue.ListItemsRequest],
) (*connect.Response[catalogue.ListItemsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_items",
			&catalogue.ListItemsContext{
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

	values := make([]*catalogue.Item, 0, len(Items))
	for _, value := range Items {
		values = append(values, value)
	}

	return connect.NewResponse(&catalogue.ListItemsResponse{Items: values}), nil
}
