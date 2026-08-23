import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
import { getUnitPrice } from "$lib/unit";
import type { Loadout as GameLoadout } from "$lib/game";

export const getLoadoutPrice = (loadout: Loadout | GameLoadout) =>
	getUnitPrice(loadout.unit!) + (loadout.item?.price || 0);
