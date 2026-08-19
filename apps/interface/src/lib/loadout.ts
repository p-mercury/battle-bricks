import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
import { getUnitPrice } from "$lib/unit";

export const getLoadoutPrice = (loadout: Loadout) =>
	getUnitPrice(loadout.unit!) + (loadout.item?.price || 0);
