package main

import (
	"connectkit"
	"context"
	"contracts/dist/catalogue/v1"
	"log/slog"

	"connectrpc.com/connect"
)

var Loadouts = map[string]*catalogue.Loadout{
	"uin075ORQu2xc53stif9": {
		Id:      "uin075ORQu2xc53stif9",
		Version: 1,
		Image:   new("/clone-scout.png"),
		Name:    "Clone Scout",
		Unit:    Units["vEPHxbpBSxzeNluR4b9U"],
		Items: []*catalogue.Item{
			Items["Kq0BO0W0rH2Ko0kyydLI"], // Blaster
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
		},
	},
	"Kq6zLj2eWy1LsbQEoC77": {
		Id:      "Kq6zLj2eWy1LsbQEoC77",
		Version: 1,
		Image:   new("/clone-sharpshooter.png"),
		Name:    "Clone Sharpshooter",
		Unit:    Units["AFFXxbwBD1aNFYbt25m7"],
		Items: []*catalogue.Item{
			Items["UG60YI7qn8DJryG9pmR7"], // Blaster Rifle
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
		},
	},
	"D9dM2ilvU1cRSFuEVFpo": {
		Id:      "D9dM2ilvU1cRSFuEVFpo",
		Version: 1,
		Image:   new("/clone-commander.png"),
		Name:    "Clone Commander",
		Unit:    Units["AFFXxbwBD1aNFYbt25m7"],
		Items: []*catalogue.Item{
			Items["AoM09AjZL8hUMhTYLIYY"], // Vibroblade
			Items["OiFq73Vk00ob7p8IuKAb"], // Hand Blasters
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
			Items["d2D9WroLYwyxAWaUpG1u"], // Blue Bolts
		},
	},
	"t7sTp7CbOWyEj1JSK2Cm": {
		Id:      "t7sTp7CbOWyEj1JSK2Cm",
		Version: 1,
		Image:   new("/fighter-tank.png"),
		Name:    "Fighter Tank",
		Unit:    Units["buMbZfN6zRflygTOVHaG"],
		Items: []*catalogue.Item{
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Blue Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
		},
	},
	"UoVeL9CfLvCgBNoTHTQK": {
		Id:      "UoVeL9CfLvCgBNoTHTQK",
		Version: 1,
		Image:   new("/droid-scout.png"),
		Name:    "Droid Scout",
		Unit:    Units["5LDT5irCFiLuKt6wQMvo"],
		Items: []*catalogue.Item{
			Items["Kq0BO0W0rH2Ko0kyydLI"], // Blaster
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"gaex5W3Mx0VbnW3yfRca": {
		Id:      "gaex5W3Mx0VbnW3yfRca",
		Version: 1,
		Image:   new("/super-battle-droid.png"),
		Name:    "Super Battle Droid",
		Unit:    Units["6EpGE8Td4gLWkcYgDIW9"],
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"kxs4sQVfzedpysp1c7x3": {
		Id:      "kxs4sQVfzedpysp1c7x3",
		Version: 1,
		Image:   new("/dwarf-spider-droid.png"),
		Name:    "Dwarf Spider Droid",
		Unit:    Units["R4ov67MKvT4YFauIDJnV"],
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
		},
	},
	"0gu35p38EmhMyLDSYAqc": {
		Id:      "0gu35p38EmhMyLDSYAqc",
		Version: 1,
		Image:   new("/octuptarra.png"),
		Name:    "Octuptarra",
		Unit:    Units["19eF4NtFTiuBXPCLFCNU"],
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
		},
	},
	"xOFvr5DoBYH2SJiCWvaK": {
		Id:      "xOFvr5DoBYH2SJiCWvaK",
		Version: 1,
		Image:   new("/armored-assault-tank.png"),
		Name:    "Armored Assault Tank",
		Unit:    Units["7JSXGhVxkJUuNjZFp3KY"],
		Items: []*catalogue.Item{
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["CV1AILTJNoyVKD7Uon63"], // Red Bolts
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["bJeepLfSKSzI1yznvTUM"], // Red Shell
			Items["xl4SlN7caFSq68DPxhr3"], // Fragmentation rocket
			Items["2A06zTJVEAHqG8Ax4w7k"], // Ion rocket
		},
	},
}

func (s *Handler) ListLoadouts(
	ctx context.Context,
	req *connect.Request[catalogue.ListLoadoutsRequest],
) (*connect.Response[catalogue.ListLoadoutsResponse], error) {
	logger := connectkit.GetLogger(ctx)
	authCtx := connectkit.GetAuthContext(ctx)

	if authCtx.Lambda != nil {
		resp, err := authCtx.Evaluate(ctx, "catalogue", "list_loadouts",
			&catalogue.ListLoadoutsContext{
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

	values := make([]*catalogue.Loadout, 0, len(Loadouts))
	for _, value := range Loadouts {
		values = append(values, value)
	}

	return connect.NewResponse(&catalogue.ListLoadoutsResponse{Loadouts: values}), nil
}
