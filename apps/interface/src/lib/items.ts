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

export interface Weapon extends Item {
	type: "WEAPON";
	ammunitionType: AmmunitionType;
	range: Range;
	fireRate: number;
}

export const items: { [key: string]: Ammunition | Weapon } = {
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
	"WEAPON#1": {
		id: "WEAPON#1",
		type: "WEAPON",
		name: "Blaster",
		price: 25,
		weight: 4,
		ammunitionType: "PLASMA",
		range: {
			min: 2,
			max: 12,
		},
		fireRate: 5,
	},
	"WEAPON#2": {
		id: "WEAPON#2",
		type: "WEAPON",
		name: "Hand Blaster",
		price: 25,
		weight: 2.5,
		ammunitionType: "PLASMA",
		range: {
			min: 0,
			max: 8,
		},
		fireRate: 4,
	},
	"WEAPON#3": {
		id: "WEAPON#3",
		type: "WEAPON",
		name: "Blaster Rifle",
		price: 40,
		weight: 5,
		ammunitionType: "PLASMA",
		range: {
			min: 4,
			max: 20,
		},
		fireRate: 3,
	},
};
