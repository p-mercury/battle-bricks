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
		name: "Droid Scout",
		unit: "UNIT#2",
		items: ["RANGE_WEAPON#1", "AMMUNITION#1", "AMMUNITION#1"],
	},
	"LOADOUT#3": {
		id: "LOADOUT#3",
		name: "Super Battle Droid ",
		unit: "UNIT#3",
		items: ["AMMUNITION#1", "AMMUNITION#1"],
	},
};
