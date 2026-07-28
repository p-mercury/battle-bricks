import type { MeleeWeapon, RangeWeapon } from "$lib/items";

export interface Unit {
	id: string;
	name: string;
	description?: string;
	price: number;
	size: number;
	speed: number;
	health: number;
	armorClass: number;
	carryCapacity: number;
	carryTypes: ("AMMUNITION" | "RANGE_WEAPON" | "MELEE_WEAPON")[];

	marksmanship: number;
	meleeAbility: number;

	items: (RangeWeapon | MeleeWeapon)[];
}

export const units: { [key: string]: Unit } = {
	"UNIT#1": {
		id: "UNIT#1",
		name: "Clone Trooper",
		price: 160,
		size: 1,
		speed: 2,
		health: 10,
		armorClass: 3,
		carryCapacity: 15,
		carryTypes: ["AMMUNITION", "RANGE_WEAPON", "MELEE_WEAPON"],
		marksmanship: 2,
		meleeAbility: 2,
		items: [
			{
				id: "",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-3",
				attackSpeed: 4,
			},
		],
	},
	"UNIT#2": {
		id: "UNIT#2",
		name: "Droid",
		price: 50,
		size: 1,
		speed: 2,
		health: 8,
		armorClass: 1,
		carryCapacity: 20,
		carryTypes: ["AMMUNITION", "RANGE_WEAPON", "MELEE_WEAPON"],
		marksmanship: 1,
		meleeAbility: 1,
		items: [
			{
				id: "",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-3",
				attackSpeed: 3,
			},
		],
	},
	"UNIT#3": {
		id: "UNIT#3",
		name: "Super Battle Droid",
		price: 150,
		size: 1,
		speed: 2,
		health: 8,
		armorClass: 2,
		carryCapacity: 10,
		carryTypes: ["AMMUNITION"],
		marksmanship: 3,
		meleeAbility: 0,
		items: [
			{
				id: "",
				type: "RANGE_WEAPON",
				name: "Arm blasters",
				price: 0,
				weight: 0,
				ammunitionType: "PLASMA",
				range: {
					min: 0,
					max: 8,
				},
				fireRate: 8,
			},
		],
	},
};
