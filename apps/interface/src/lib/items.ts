export interface Item {
	id: string;
	type: string;
	name: string;
	description?: string;
	price: number;
	weight: number;
}

export type AmmunitionType = "PLASMA" | "SLUG";

export interface Range {
	min: number;
	max: number;
}

export interface Ammunition extends Item {
	type: "AMMUNITION";
	ammunitionType: AmmunitionType;
	capacity: number;
	damage: string;
	armorPiercing: number;
}

export interface RangeWeapon extends Item {
	type: "RANGE_WEAPON";
	ammunitionType: AmmunitionType;
	range: Range;
	fireRate: number;
}

export interface MeleeWeapon extends Item {
	type: "MELEE_WEAPON";
	armorPiercing: number;
	damage: string;
	attackSpeed: number;
}

export const items: { [key: string]: Ammunition | RangeWeapon | MeleeWeapon } =
	{
		"AMMUNITION#1": {
			id: "AMMUNITION#1",
			type: "AMMUNITION",
			name: "Red Plasma Cartridge",
			price: 5,
			weight: 3,
			ammunitionType: "PLASMA",
			capacity: 160,
			damage: "1d6-1",
			armorPiercing: 1,
		},
		"AMMUNITION#2": {
			id: "AMMUNITION#2",
			type: "AMMUNITION",
			name: "Blue Plasma Cartridge",
			price: 10,
			weight: 3,
			ammunitionType: "PLASMA",
			capacity: 100,
			damage: "1d6",
			armorPiercing: 2,
		},
		"AMMUNITION#3": {
			id: "AMMUNITION#3",
			type: "AMMUNITION",
			name: "Green Plasma Cartridge",
			price: 10,
			weight: 3,
			ammunitionType: "PLASMA",
			capacity: 100,
			damage: "1d6+1",
			armorPiercing: 2,
		},
		"AMMUNITION#4": {
			id: "AMMUNITION#4",
			type: "AMMUNITION",
			name: "Yellow Plasma Cartridge",
			price: 10,
			weight: 3,
			ammunitionType: "PLASMA",
			capacity: 80,
			damage: "1d6+1",
			armorPiercing: 3,
		},
		"RANGE_WEAPON#1": {
			id: "RANGE_WEAPON#1",
			type: "RANGE_WEAPON",
			name: "Blaster",
			price: 20,
			weight: 4,
			ammunitionType: "PLASMA",
			range: {
				min: 2,
				max: 12,
			},
			fireRate: 5,
		},
		"RANGE_WEAPON#2": {
			id: "RANGE_WEAPON#2",
			type: "RANGE_WEAPON",
			name: "Hand Blaster",
			price: 20,
			weight: 2.5,
			ammunitionType: "PLASMA",
			range: {
				min: 0,
				max: 8,
			},
			fireRate: 4,
		},
		"RANGE_WEAPON#3": {
			id: "RANGE_WEAPON#3",
			type: "RANGE_WEAPON",
			name: "Blaster Rifle",
			price: 35,
			weight: 5,
			ammunitionType: "PLASMA",
			range: {
				min: 4,
				max: 20,
			},
			fireRate: 3,
		},
		"MELEE_WEAPON#1": {
			id: "MELEE_WEAPON#1",
			type: "MELEE_WEAPON",
			name: "Vibroblade",
			price: 15,
			weight: 1,
			armorPiercing: 1,
			damage: "1d6-1",
			attackSpeed: 5,
		},
		"MELEE_WEAPON#2": {
			id: "MELEE_WEAPON#2",
			type: "MELEE_WEAPON",
			name: "Vibrosword",
			price: 30,
			weight: 2,
			armorPiercing: 2,
			damage: "1d6",
			attackSpeed: 4,
		},
	};
