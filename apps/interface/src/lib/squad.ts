import type { Faction } from "$lib/units";

export interface Squad {
	id: string;
	name: string;
	faction: Faction;
	loadouts: string[];
}
