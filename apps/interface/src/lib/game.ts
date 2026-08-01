import type { Faction } from "$lib/faction";
import type { Loadout } from "$lib/loadouts";
import type { Unit } from "$lib/units";

export interface Game {
	initiative: number;
	attacker: {
		name: string;
		faction: Faction;
		loadouts: (Loadout & { initiative: number })[];
	};
	defender: {
		faction: Faction;
		units: Unit[];
	};
}
