import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";

export const getLoadoutPrice = (loadout: Loadout) =>
	loadout.unit!.price + loadout.items.reduce((b, item) => b + item.price, 0);
