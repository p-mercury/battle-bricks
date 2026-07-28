import type { Loadout } from "$lib/loadouts";
import { units } from "$lib/units";
import { items } from "$lib/items";

export function getLoadoutStats(loadout: Loadout) {
	const lUnit = units[loadout.unit];
	let lPrice = lUnit.price;
	let lCarryWeight = 0;
	const lItems = [
		...lUnit.items,
		...loadout.items.map((i) => {
			const item = items[i];
			lPrice += item.price;
			lCarryWeight += item.weight;
			return item;
		}),
	];

	return {
		unit: lUnit,
		items: lItems,
		price: lPrice,
		carryWeight: lCarryWeight,
	};
}
