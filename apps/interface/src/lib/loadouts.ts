export interface Loadout {
	id: string;
	name: string;
	description?: string;
	unit: string;
	items: string[];
}

export const loadouts: { [key: string]: Loadout } = {
	"LOADOUT#1": {
		id: "LOADOUT#1",
		name: "Clone Trooper",
		unit: "UNIT#1",
		items: ["RANGE_WEAPON#1", "AMMUNITION#2", "AMMUNITION#2"],
	},
	"LOADOUT#2": {
		id: "LOADOUT#2",
		name: "Clone Sharpshooter",
		unit: "UNIT#2",
		items: ["RANGE_WEAPON#3", "AMMUNITION#2"],
	},
	"LOADOUT#3": {
		id: "LOADOUT#3",
		name: "Fighter Tank",
		unit: "UNIT#3",
		items: ["AMMUNITION#2", "AMMUNITION#2", "AMMUNITION#5", "AMMUNITION#5"],
	},
	"LOADOUT#4": {
		id: "LOADOUT#4",
		name: "Droid Scout",
		unit: "UNIT#4",
		items: ["RANGE_WEAPON#1", "AMMUNITION#1", "AMMUNITION#1"],
	},
	"LOADOUT#5": {
		id: "LOADOUT#5",
		name: "Super Battle Droid",
		unit: "UNIT#5",
		items: ["AMMUNITION#1", "AMMUNITION#1"],
	},
	"LOADOUT#6": {
		id: "LOADOUT#6",
		name: "Octuptarra",
		unit: "UNIT#6",
		items: [
			"AMMUNITION#1",
			"AMMUNITION#1",
			"AMMUNITION#5",
			"AMMUNITION#5",
			"AMMUNITION#5",
		],
	},
};
