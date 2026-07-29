import type { MeleeWeapon, RangeWeapon } from "$lib/items";

export type Faction =
	| "GALACTIC_REPUBLIC"
	| "REBEL_ALLIANCE"
	| "SEPARATIST_ALLIANCE"
	| "GALACTIC_EMPIRE";

export interface Unit {
	id: string;
	faction: Faction[];
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
		faction: ["GALACTIC_REPUBLIC"],
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
		faction: ["GALACTIC_REPUBLIC"],
		name: "Clone Specialist",
		price: 200,
		size: 1,
		speed: 2,
		health: 10,
		armorClass: 3,
		carryCapacity: 15,
		carryTypes: ["AMMUNITION", "RANGE_WEAPON", "MELEE_WEAPON"],
		marksmanship: 3,
		meleeAbility: 2,
		items: [
			{
				id: "",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-1",
				attackSpeed: 4,
			},
		],
	},
	"UNIT#3": {
		id: "UNIT#3",
		faction: ["SEPARATIST_ALLIANCE"],
		name: "Droid",
		price: 50,
		size: 1,
		speed: 2,
		health: 8,
		armorClass: 1,
		carryCapacity: 15,
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
	"UNIT#4": {
		id: "UNIT#4",
		faction: ["SEPARATIST_ALLIANCE"],
		name: "Super Battle Droid",
		price: 150,
		size: 1,
		speed: 2,
		health: 8,
		armorClass: 2,
		carryCapacity: 9,
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
	"UNIT#5": {
		id: "UNIT#5",
		faction: ["SEPARATIST_ALLIANCE"],
		name: "Octuptarra",
		price: 350,
		size: 3,
		speed: 3,
		health: 22,
		armorClass: 3,
		carryCapacity: 24,
		carryTypes: ["AMMUNITION"],
		marksmanship: 3,
		meleeAbility: 0,
		items: [
			{
				id: "",
				type: "RANGE_WEAPON",
				name: "Laser Cannon",
				price: 0,
				weight: 0,
				ammunitionType: "PLASMA",
				range: {
					min: 4,
					max: 10,
				},
				fireRate: 6,
			},
			{
				id: "",
				type: "RANGE_WEAPON",
				name: "Rocket Launcher",
				price: 0,
				weight: 0,
				ammunitionType: "ROCKET",
				range: {
					min: 5,
					max: 12,
				},
				fireRate: 1,
			},
		],
	},
};
