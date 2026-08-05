import { units, type Unit } from "$lib/units";
import {
	items,
	type BoltAmmunition,
	type ShellAmmunition,
	type RocketAmmunition,
	type MeleeWeapon,
	type RangeWeapon,
} from "$lib/items";

export interface LoadoutConfig {
	id: string;
	image?: string;
	name: string;
	description?: string;
	unit: string;
	items: string[];
}

const loadoutConfigs: { [key: string]: LoadoutConfig } = {
	"LOADOUT#1": {
		id: "LOADOUT#1",
		image: "/clone-scout.png",
		name: "Clone Scout",
		unit: "UNIT#1",
		items: ["RANGE_WEAPON#1", "BOLT_AMMUNITION#2", "BOLT_AMMUNITION#2"],
	},
	"LOADOUT#2": {
		id: "LOADOUT#2",
		image: "/clone-sharpshooter.png",
		name: "Clone Sharpshooter",
		unit: "UNIT#2",
		items: ["RANGE_WEAPON#3", "BOLT_AMMUNITION#2"],
	},

	"LOADOUT#3": {
		id: "LOADOUT#3",
		image: "/clone-commander.png",
		name: "Clone Commander",
		unit: "UNIT#2",
		items: [
			"MELEE_WEAPON#1",
			"RANGE_WEAPON#2",
			"BOLT_AMMUNITION#2",
			"BOLT_AMMUNITION#2",
		],
	},
	"LOADOUT#4": {
		id: "LOADOUT#4",
		name: "Fighter Tank",
		image: "/fighter-tank.png",
		unit: "UNIT#3",
		items: [
			"SHELL_AMMUNITION#2",
			"SHELL_AMMUNITION#2",
			"SHELL_AMMUNITION#2",
			"SHELL_AMMUNITION#2",
			"ROCKET_AMMUNITION#1",
			"ROCKET_AMMUNITION#2",
		],
	},
	"LOADOUT#5": {
		id: "LOADOUT#5",
		image: "/droid-scout.png",
		name: "Droid Scout",
		unit: "UNIT#4",
		items: ["RANGE_WEAPON#1", "BOLT_AMMUNITION#1"],
	},
	"LOADOUT#6": {
		id: "LOADOUT#6",
		image: "/super-battle-droid.png",
		name: "Super Battle Droid",
		unit: "UNIT#5",
		items: ["BOLT_AMMUNITION#1", "BOLT_AMMUNITION#1"],
	},
	"LOADOUT#7": {
		id: "LOADOUT#7",
		image: "/octuptarra.png",
		name: "Octuptarra",
		unit: "UNIT#6",
		items: [
			"BOLT_AMMUNITION#1",
			"BOLT_AMMUNITION#1",
			"ROCKET_AMMUNITION#1",
			"ROCKET_AMMUNITION#1",
			"ROCKET_AMMUNITION#2",
		],
	},
	"LOADOUT#8": {
		id: "LOADOUT#8",
		image: "/armored-assault-tank.png",
		name: "Armored Assault Tank",
		unit: "UNIT#7",
		items: [
			"BOLT_AMMUNITION#1",
			"BOLT_AMMUNITION#1",
			"SHELL_AMMUNITION#1",
			"SHELL_AMMUNITION#1",
			"SHELL_AMMUNITION#1",
			"ROCKET_AMMUNITION#1",
			"ROCKET_AMMUNITION#2",
		],
	},
};

export interface Loadout {
	id: string;
	image?: string;
	name: string;
	unit: Unit;
	items: (
		| BoltAmmunition
		| ShellAmmunition
		| RocketAmmunition
		| RangeWeapon
		| MeleeWeapon
	)[];
	price: number;
	carryWeight: number;
}

export const loadouts = Object.fromEntries(
	Object.entries(loadoutConfigs).map(([key, loadout]) => {
		const lUnit = units[loadout.unit];
		let lPrice = lUnit.price;
		let lCarryWeight = 0;
		const lItems = [
			...loadout.items.map((i) => {
				const item = items[i];
				lPrice += item.price;
				lCarryWeight += item.weight;
				return item;
			}),
		];

		return [
			key,
			{
				id: loadout.id,
				image: loadout.image,
				name: loadout.name,
				unit: lUnit,
				items: lItems,
				price: lPrice,
				carryWeight: lCarryWeight,
			},
		] as [string, Loadout];
	}),
) as Record<string, Loadout>;
