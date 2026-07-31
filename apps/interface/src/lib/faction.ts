export type Faction =
	| "GALACTIC_REPUBLIC"
	| "REBEL_ALLIANCE"
	| "SEPARATIST_ALLIANCE"
	| "GALACTIC_EMPIRE";

export function getFactionName(t: Faction): string {
	switch (t) {
		case "GALACTIC_REPUBLIC":
			return /* @wc-ignore */ "Galactic Republic";
		case "REBEL_ALLIANCE":
			return /* @wc-ignore */ "Rebel Alliance";
		case "SEPARATIST_ALLIANCE":
			return /* @wc-ignore */ "Separatist Alliance";
		case "GALACTIC_EMPIRE":
			return /* @wc-ignore */ "Galactic Empire";
		default:
			return /* @wc-include */ "Unknown";
	}
}
