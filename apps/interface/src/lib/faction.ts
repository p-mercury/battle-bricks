export type Faction =
	| "GALACTIC_REPUBLIC"
	| "REBEL_ALLIANCE"
	| "SEPARATIST_ALLIANCE"
	| "GALACTIC_EMPIRE";

export function getFactionName(t: Faction): string {
	switch (t) {
		case "GALACTIC_REPUBLIC":
			return "Galactic Republic";
		case "REBEL_ALLIANCE":
			return "Rebel Alliance";
		case "SEPARATIST_ALLIANCE":
			return "Separatist Alliance";
		case "GALACTIC_EMPIRE":
			return "Galactic Empire";
		default:
			return "Unknown";
	}
}
