import type { Faction } from "$lib/faction";
import type { Loadout } from "$lib/loadouts";
import type { Unit } from "$lib/units";

export interface GameLoadout extends Loadout {
	color: string;
	turnComplete: boolean;
	inCover?: boolean;
}

export interface Game {
	attacker: {
		name: string;
		faction: Faction;
		loadouts: GameLoadout[];
	};
	defender: {
		faction: Faction;
		units: Unit[];
	};
}
