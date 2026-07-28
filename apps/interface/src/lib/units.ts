export interface Unit {
	id: string;
	name: string;
	description?: string;
	price: number;
	size: number;
	health: number;
	armorClass: number;
	carryCapacity: number;
	accuracy: number;
	items: string[];
}

export const units: { [key: string]: Unit } = {
	"UNIT#1": {
		id: "UNIT#1",
		name: "Clone Trooper",
		price: 160,
		size: 1,
		health: 10,
		armorClass: 3,
		carryCapacity: 15,
		accuracy: 2,
		items: [],
	},
	"UNIT#2": {
		id: "UNIT#2",
		name: "Droid",
		price: 50,
		size: 1,
		health: 8,
		armorClass: 1,
		carryCapacity: 20,
		accuracy: 1,
		items: [],
	},
	"UNIT#3": {
		id: "UNIT#3",
		name: "Super Battle Droid",
		price: 150,
		size: 1,
		health: 8,
		armorClass: 2,
		carryCapacity: 20,
		accuracy: 2,
		items: ["WEAPON#1"],
	},
};
