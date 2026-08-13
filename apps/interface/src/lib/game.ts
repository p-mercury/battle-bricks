import type { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";

export interface GameLoadout extends Loadout {
	color: string;
	turnComplete: boolean;
	inCover?: boolean;
}

export interface Game {
	attacker: {
		name: string;
		faction: Faction;
		loadouts: { [key: string]: GameLoadout };
	};
	defender: {
		faction: Faction;
		units: Unit[];
	};
}
