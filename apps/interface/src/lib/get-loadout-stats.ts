import type { Loadout } from "$lib/loadouts";
import { units } from "$lib/units";
import { items } from "$lib/items";

export function getLoadoutStats(loadout: Loadout) {
	const lUnit = units[loadout.unit];
	let lPrice = 0;
	let lCarryWeight = 0;
	const lItems = [...loadout.items, ...lUnit.items].map((i) => {
		const item = items[i];
		lPrice += item.price;
		lCarryWeight += item.weight;
		return item;
	});

	return {
		unit: lUnit,
		items: lItems,
		price: lPrice,
		carryWeight: lCarryWeight,
	};
}
