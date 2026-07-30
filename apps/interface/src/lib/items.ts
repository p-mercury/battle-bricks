export interface Item {
	id: string;
	type: string;
	name: string;
	description?: string;
	price: number;
	weight: number;
}

export type AmmunitionType = "PLASMA" | "ROCKET" | "SLUG";

export interface Range {
	min: number;
	max: number;
}

export interface PlasmaAmmunition extends Item {
	type: "PLASMA_AMMUNITION";
	ammunitionType: "PLASMA";
	capacity: number;
	damage: string;
	armorPiercing: number;
	splashRadius: number;
}

export interface RocketAmmunition extends Item {
	type: "ROCKET_AMMUNITION";
	ammunitionType: "ROCKET";
	damage: string;
	splashRadius: number;
}

export interface SlugAmmunition extends Item {
	type: "SLUG_AMMUNITION";
	ammunitionType: "SLUG";
	capacity: number;
	damage: string;
	armorPiercing: number;
	splashRadius: number;
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

export const items: {
	[key: string]:
		| PlasmaAmmunition
		| RocketAmmunition
		| SlugAmmunition
		| RangeWeapon
		| MeleeWeapon;
} = {
	"PLASMA_AMMUNITION#1": {
		id: "PLASMA_AMMUNITION#1",
		type: "PLASMA_AMMUNITION",
		name: "Red Plasma Cartridge",
		price: 5,
		weight: 3,
		ammunitionType: "PLASMA",
		capacity: 100,
		damage: "1d6-1",
		armorPiercing: 1,
		splashRadius: 0,
	},
	"PLASMA_AMMUNITION#2": {
		id: "PLASMA_AMMUNITION#2",
		type: "PLASMA_AMMUNITION",
		name: "Blue Plasma Cartridge",
		price: 10,
		weight: 3,
		ammunitionType: "PLASMA",
		capacity: 60,
		damage: "1d6",
		armorPiercing: 2,
		splashRadius: 0,
	},
	"PLASMA_AMMUNITION#3": {
		id: "PLASMA_AMMUNITION#3",
		type: "PLASMA_AMMUNITION",
		name: "Green Plasma Cartridge",
		price: 15,
		weight: 3,
		ammunitionType: "PLASMA",
		capacity: 60,
		damage: "1d6+1",
		armorPiercing: 2,
		splashRadius: 0,
	},
	"PLASMA_AMMUNITION#4": {
		id: "PLASMA_AMMUNITIONN#4",
		type: "PLASMA_AMMUNITION",
		name: "Yellow Plasma Cartridge",
		price: 20,
		weight: 3,
		ammunitionType: "PLASMA",
		capacity: 50,
		damage: "1d6+1",
		armorPiercing: 3,
		splashRadius: 0,
	},
	"ROCKET_AMMUNITION#1": {
		id: "ROCKET_AMMUNITION#1",
		type: "ROCKET_AMMUNITION",
		name: "Fragmentation rocket",
		price: 30,
		weight: 5,
		ammunitionType: "ROCKET",
		damage: "1d6+2",
		splashRadius: 2,
	},
	"ROCKET_AMMUNITION#2": {
		id: "ROCKET_AMMUNITION#2",
		type: "ROCKET_AMMUNITION",
		name: "Ion rocket",
		price: 40,
		weight: 5,
		ammunitionType: "ROCKET",
		damage: "1d6+4",
		splashRadius: 1,
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
		fireRate: 6,
	},
	"RANGE_WEAPON#2": {
		id: "RANGE_WEAPON#2",
		type: "RANGE_WEAPON",
		name: "Hand Blasters",
		price: 40,
		weight: 5,
		ammunitionType: "PLASMA",
		range: {
			min: 0,
			max: 8,
		},
		fireRate: 8,
	},
	"RANGE_WEAPON#3": {
		id: "RANGE_WEAPON#3",
		type: "RANGE_WEAPON",
		name: "Blaster Rifle",
		price: 40,
		weight: 5,
		ammunitionType: "PLASMA",
		range: {
			min: 5,
			max: 20,
		},
		fireRate: 4,
	},
	"MELEE_WEAPON#1": {
		id: "MELEE_WEAPON#1",
		type: "MELEE_WEAPON",
		name: "Vibroblade",
		price: 20,
		weight: 1,
		armorPiercing: 2,
		damage: "1d6",
		attackSpeed: 5,
	},
	"MELEE_WEAPON#2": {
		id: "MELEE_WEAPON#2",
		type: "MELEE_WEAPON",
		name: "Vibrosword",
		price: 30,
		weight: 2,
		armorPiercing: 2,
		damage: "1d6+1",
		attackSpeed: 4,
	},
};
