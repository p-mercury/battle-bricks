export interface Item {
	id: string;
	type: string;
	name: string;
	description?: string;
	price: number;
	weight: number;
}

export type AmmunitionType = "BOLT" | "SHELL" | "ROCKET";

export interface Range {
	min: number;
	max: number;
}

export interface BoltAmmunition extends Item {
	type: "BOLT_AMMUNITION";
	ammunitionType: "BOLT";
	capacity: number;
	damage: string;
	armorPiercing: number;
}

export interface ShellAmmunition extends Item {
	type: "SHELL_AMMUNITION";
	ammunitionType: "SHELL";
	damage: string;
	armorPiercing: number;
}

export interface RocketAmmunition extends Item {
	type: "ROCKET_AMMUNITION";
	ammunitionType: "ROCKET";
	damage: string;
	splashRadius: number;
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

export interface Action extends Item {
	type: "ACTION";
}

export const items: {
	[key: string]:
		| BoltAmmunition
		| ShellAmmunition
		| RocketAmmunition
		| RangeWeapon
		| MeleeWeapon
		| Action;
} = {
	"BOLT_AMMUNITION#1": {
		id: "BOLT_AMMUNITION#1",
		type: "BOLT_AMMUNITION",
		name: "Red Plasma Cartridge",
		description: "Cheap Tibanna gas blend producing weaker red bolts",
		price: 5,
		weight: 2,
		ammunitionType: "BOLT",
		capacity: 120,
		damage: "1d6-1",
		armorPiercing: 1,
	},
	"BOLT_AMMUNITION#2": {
		id: "BOLT_AMMUNITION#2",
		type: "BOLT_AMMUNITION",
		name: "Blue Plasma Cartridge",
		description: "High grade Tibanna gas blend producing blue bolts",
		price: 10,
		weight: 2,
		ammunitionType: "BOLT",
		capacity: 90,
		damage: "1d6",
		armorPiercing: 2,
	},
	"BOLT_AMMUNITION#3": {
		id: "BOLT_AMMUNITION#3",
		type: "BOLT_AMMUNITION",
		name: "Green Plasma Cartridge",
		description: "Pure refined Tibanna producing powerful green bolts",
		price: 15,
		weight: 2,
		ammunitionType: "BOLT",
		capacity: 90,
		damage: "1d6+1",
		armorPiercing: 2,
	},
	"BOLT_AMMUNITION#4": {
		id: "BOLT_AMMUNITIONN#4",
		type: "BOLT_AMMUNITION",
		name: "Yellow Plasma Cartridge",
		description: "High pressure Tibanna producing armor piercing yellow bolts",
		price: 20,
		weight: 2,
		ammunitionType: "BOLT",
		capacity: 60,
		damage: "1d6+1",
		armorPiercing: 3,
	},
	"SHELL_AMMUNITION#1": {
		id: "SHELL_AMMUNITION#1",
		type: "SHELL_AMMUNITION",
		name: "Red Plasma Shell",
		description: "Cheap Tibanna gas blend producing weaker blasts",
		price: 10,
		weight: 4,
		ammunitionType: "SHELL",
		damage: "2d6",
		armorPiercing: 3,
	},
	"SHELL_AMMUNITION#2": {
		id: "SHELL_AMMUNITION#2",
		type: "SHELL_AMMUNITION",
		name: "Blue Plasma Shell",
		description: "High grade Tibanna gas blend producing blue bolts",
		price: 20,
		weight: 4,
		ammunitionType: "SHELL",
		damage: "2d6+2",
		armorPiercing: 3,
	},
	"SHELL_AMMUNITION#3": {
		id: "SHELL_AMMUNITION#3",
		type: "SHELL_AMMUNITION",
		name: "Green Plasma Shell",
		description: "Pure refined Tibanna producing powerful green bolts",
		price: 30,
		weight: 4,
		ammunitionType: "SHELL",
		damage: "2d6+4",
		armorPiercing: 4,
	},
	"ROCKET_AMMUNITION#1": {
		id: "ROCKET_AMMUNITION#1",
		type: "ROCKET_AMMUNITION",
		name: "Fragmentation rocket",
		description:
			"Anti-personnel rocket that shreds unarmored targets with a wide blast radius",
		price: 30,
		weight: 6,
		ammunitionType: "ROCKET",
		damage: "1d6+4",
		armorPiercing: 1,
		splashRadius: 2,
	},
	"ROCKET_AMMUNITION#2": {
		id: "ROCKET_AMMUNITION#2",
		type: "ROCKET_AMMUNITION",
		name: "Ion rocket",
		price: 30,
		weight: 6,
		ammunitionType: "ROCKET",
		damage: "1d6+4",
		armorPiercing: 4,
		splashRadius: 0,
	},
	"RANGE_WEAPON#1": {
		id: "RANGE_WEAPON#1",
		type: "RANGE_WEAPON",
		name: "Blaster",
		price: 20,
		weight: 4,
		ammunitionType: "BOLT",
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
		ammunitionType: "BOLT",
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
		price: 30,
		weight: 5,
		ammunitionType: "BOLT",
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
