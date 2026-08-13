import type { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";

export interface Squad {
	id: string;
	name: string;
	faction: Faction;
	loadouts: string[];
}
