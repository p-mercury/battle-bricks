import type { MeleeWeapon, RangeWeapon, Action } from "$lib/items";
import type { Faction } from "$lib/faction";

export interface Unit {
	id: string;
	faction: Faction[];
	image?: string;
	name: string;
	description?: string;
	price: number;
	size: number;
	speed: number;
	health: number;
	armorClass: number;
	carryCapacity: number;

	marksmanship?: number;
	meleeAbility?: number;

	items: (RangeWeapon | MeleeWeapon | Action)[];
}

export const units: { [key: string]: Unit } = {
	"UNIT#1": {
		id: "UNIT#1",
		faction: ["GALACTIC_REPUBLIC"],
		image: "/unit-clone.png",
		name: "Clone Trooper",
		price: 150,
		size: 1,
		speed: 2,
		health: 10,
		armorClass: 3,
		carryCapacity: 15,
		marksmanship: 2,
		meleeAbility: 2,
		items: [
			{
				id: "INTERNAL#1",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-3",
				attackSpeed: 4,
			},
			{
				id: "INTERNAL#2",
				type: "ACTION",
				name: "Sprint",
				price: 0,
				weight: 0,
				description: "Double move",
			},
		],
	},
	"UNIT#2": {
		id: "UNIT#2",
		faction: ["GALACTIC_REPUBLIC"],
		image: "/unit-clone-specialist.png",
		name: "Clone Specialist",
		price: 190,
		size: 1,
		speed: 2,
		health: 10,
		armorClass: 3,
		carryCapacity: 15,
		marksmanship: 3,
		meleeAbility: 2,
		items: [
			{
				id: "INTERNAL#1",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-1",
				attackSpeed: 4,
			},
			{
				id: "INTERNAL#2",
				type: "ACTION",
				name: "Sprint",
				price: 0,
				weight: 0,
				description: "Double move",
			},
		],
	},
	"UNIT#3": {
		id: "UNIT#3",
		faction: ["GALACTIC_REPUBLIC"],
		image: "/fighter-tank.png",
		name: "Fighter Tank",
		price: 350,
		size: 4,
		speed: 4,
		health: 26,
		armorClass: 4,
		carryCapacity: 24,
		marksmanship: 3,
		items: [
			{
				id: "INTERNAL#1",
				type: "RANGE_WEAPON",
				name: "Laser Cannons",
				price: 0,
				weight: 0,
				ammunitionType: "SHELL",
				range: {
					min: 5,
					max: 12,
				},
				fireRate: 2,
			},
			{
				id: "INTERNAL#2",
				type: "RANGE_WEAPON",
				name: "Rocket Launcher",
				price: 0,
				weight: 0,
				ammunitionType: "ROCKET",
				range: {
					min: 6,
					max: 12,
				},
				fireRate: 2,
			},
		],
	},
	"UNIT#4": {
		id: "UNIT#4",
		faction: ["SEPARATIST_ALLIANCE"],
		image: "/unit-droid.png",
		name: "Droid",
		price: 55,
		size: 1,
		speed: 2,
		health: 8,
		armorClass: 1,
		carryCapacity: 15,
		marksmanship: 1,
		meleeAbility: 1,
		items: [
			{
				id: "INTERNAL#1",
				type: "MELEE_WEAPON",
				name: "Unarmed Strike",
				price: 0,
				weight: 0,
				armorPiercing: 1,
				damage: "1d6-3",
				attackSpeed: 3,
			},
			{
				id: "INTERNAL#2",
				type: "ACTION",
				name: "Sprint",
				price: 0,
				weight: 0,
				description: "Double move",
			},
		],
	},
	"UNIT#5": {
		id: "UNIT#5",
		faction: ["SEPARATIST_ALLIANCE"],
		image: "/unit-super-battle-droid.png",
		name: "Super Battle Droid",
		price: 160,
		size: 1,
		speed: 1,
		health: 10,
		armorClass: 2,
		carryCapacity: 9,
		marksmanship: 2,
		items: [
			{
				id: "INTERNAL#1",
				type: "RANGE_WEAPON",
				name: "Arm blasters",
				price: 0,
				weight: 0,
				ammunitionType: "BOLT",
				range: {
					min: 0,
					max: 8,
				},
				fireRate: 8,
			},
			{
				id: "INTERNAL#2",
				type: "ACTION",
				name: "Sprint",
				price: 0,
				weight: 0,
				description: "Double move",
			},
		],
	},
	"UNIT#6": {
		id: "UNIT#6",
		faction: ["SEPARATIST_ALLIANCE"],
		image: "/dwarf-spider-droid.png",
		name: "Dwarf Spider Droid",
		price: 360,
		size: 2,
		speed: 3,
		health: 16,
		armorClass: 2,
		carryCapacity: 24,
		marksmanship: 3,
		items: [
			{
				id: "INTERNAL#1",
				type: "RANGE_WEAPON",
				name: "Laser Blaster",
				price: 0,
				weight: 0,
				ammunitionType: "BOLT",
				range: {
					min: 2,
					max: 10,
				},
				fireRate: 8,
			},
			{
				id: "INTERNAL#2",
				type: "ACTION",
				name: "Sprint",
				price: 0,
				weight: 0,
				description: "Double move",
			},
		],
	},
	"UNIT#7": {
		id: "UNIT#7",
		faction: ["SEPARATIST_ALLIANCE"],
		image: "/octuptarra.png",
		name: "Octuptarra",
		price: 360,
		size: 4,
		speed: 3,
		health: 18,
		armorClass: 3,
		carryCapacity: 24,
		marksmanship: 3,
		items: [
			{
				id: "INTERNAL#1",
				type: "RANGE_WEAPON",
				name: "Laser Blaster",
				price: 0,
				weight: 0,
				ammunitionType: "BOLT",
				range: {
					min: 4,
					max: 10,
				},
				fireRate: 8,
			},
			{
				id: "INTERNAL#2",
				type: "RANGE_WEAPON",
				name: "Rocket Launcher",
				price: 0,
				weight: 0,
				ammunitionType: "ROCKET",
				range: {
					min: 6,
					max: 12,
				},
				fireRate: 1,
			},
		],
	},
	"UNIT#8": {
		id: "UNIT#8",
		faction: ["SEPARATIST_ALLIANCE"],
		image: "/armored-assault-tank.png",
		name: "Armored Assault Tank",
		price: 400,
		size: 4,
		speed: 3,
		health: 22,
		armorClass: 4,
		carryCapacity: 24,
		marksmanship: 2,
		items: [
			{
				id: "INTERNAL#1",
				type: "RANGE_WEAPON",
				name: "Side Blasters",
				price: 0,
				weight: 0,
				ammunitionType: "BOLT",
				range: {
					min: 4,
					max: 10,
				},
				fireRate: 8,
			},
			{
				id: "INTERNAL#2",
				type: "RANGE_WEAPON",
				name: "Main Turret",
				price: 0,
				weight: 0,
				ammunitionType: "SHELL",
				range: {
					min: 5,
					max: 12,
				},
				fireRate: 1,
			},
			{
				id: "INTERNAL#3",
				type: "RANGE_WEAPON",
				name: "Rocket Launchers",
				price: 0,
				weight: 0,
				ammunitionType: "ROCKET",
				range: {
					min: 6,
					max: 12,
				},
				fireRate: 2,
			},
		],
	},
};
