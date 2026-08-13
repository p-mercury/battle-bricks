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
		unit: "vEPHxbpBSxzeNluR4b9U",
		items: [
			"Kq0BO0W0rH2Ko0kyydLI", // Blaster
			"d2D9WroLYwyxAWaUpG1u", // Blue Bolts
			"d2D9WroLYwyxAWaUpG1u", // Blue Bolts
		],
	},
	"LOADOUT#2": {
		id: "LOADOUT#2",
		image: "/clone-sharpshooter.png",
		name: "Clone Sharpshooter",
		unit: "AFFXxbwBD1aNFYbt25m7",
		items: [
			"UG60YI7qn8DJryG9pmR7", // Blaster Rifle
			"d2D9WroLYwyxAWaUpG1u", // Blue Bolts
		],
	},

	"LOADOUT#3": {
		id: "LOADOUT#3",
		image: "/clone-commander.png",
		name: "Clone Commander",
		unit: "AFFXxbwBD1aNFYbt25m7",
		items: [
			"AoM09AjZL8hUMhTYLIYY", // Vibroblade
			"OiFq73Vk00ob7p8IuKAb", // Hand Blasters
			"d2D9WroLYwyxAWaUpG1u", // Blue Bolts
			"d2D9WroLYwyxAWaUpG1u", // Blue Bolts
		],
	},
	"LOADOUT#4": {
		id: "LOADOUT#4",
		name: "Fighter Tank",
		image: "/fighter-tank.png",
		unit: "buMbZfN6zRflygTOVHaG",
		items: [
			"bJeepLfSKSzI1yznvTUM", // Blue Shell
			"bJeepLfSKSzI1yznvTUM", // Blue Shell
			"bJeepLfSKSzI1yznvTUM", // Blue Shell
			"bJeepLfSKSzI1yznvTUM", // Blue Shell
			"xl4SlN7caFSq68DPxhr3", // Fragmentation rocket
			"2A06zTJVEAHqG8Ax4w7k", // Ion rocket
		],
	},
	"LOADOUT#5": {
		id: "LOADOUT#5",
		image: "/droid-scout.png",
		name: "Droid Scout",
		unit: "5LDT5irCFiLuKt6wQMvo",
		items: [
			"Kq0BO0W0rH2Ko0kyydLI", // Blaster
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
		],
	},
	"LOADOUT#6": {
		id: "LOADOUT#6",
		image: "/super-battle-droid.png",
		name: "Super Battle Droid",
		unit: "6EpGE8Td4gLWkcYgDIW9",
		items: [
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
		],
	},
	"LOADOUT#7": {
		id: "LOADOUT#7",
		image: "/dwarf-spider-droid.png",
		name: "Dwarf Spider Droid",
		unit: "R4ov67MKvT4YFauIDJnV",
		items: [
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
		],
	},
	"LOADOUT#8": {
		id: "LOADOUT#8",
		image: "/octuptarra.png",
		name: "Octuptarra",
		unit: "19eF4NtFTiuBXPCLFCNU",
		items: [
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"xl4SlN7caFSq68DPxhr3", // Fragmentation rocket
			"xl4SlN7caFSq68DPxhr3", // Fragmentation rocket
			"2A06zTJVEAHqG8Ax4w7k", // Ion rocket
		],
	},
	"LOADOUT#9": {
		id: "LOADOUT#9",
		image: "/armored-assault-tank.png",
		name: "Armored Assault Tank",
		unit: "7JSXGhVxkJUuNjZFp3KY",
		items: [
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"CV1AILTJNoyVKD7Uon63", // Red Bolts
			"bJeepLfSKSzI1yznvTUM", // Red Shell
			"bJeepLfSKSzI1yznvTUM", // Red Shell
			"bJeepLfSKSzI1yznvTUM", // Red Shell
			"xl4SlN7caFSq68DPxhr3", // Fragmentation rocket
			"2A06zTJVEAHqG8Ax4w7k", // Ion rocket
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
