import type { Faction } from "$lib/faction";

export interface Squad {
	id: string;
	name: string;
	faction: Faction;
	loadouts: string[];
}
