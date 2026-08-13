import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";

export function getFactionName(t: Faction): string {
	switch (t) {
		case Faction.GALACTIC_REPUBLIC:
			return "Galactic Republic";
		case Faction.REBEL_ALLIANCE:
			return "Rebel Alliance";
		case Faction.SEPARATIST_ALLIANCE:
			return "Separatist Alliance";
		case Faction.GALACTIC_EMPIRE:
			return "Galactic Empire";
		default:
			return "Unknown";
	}
}
