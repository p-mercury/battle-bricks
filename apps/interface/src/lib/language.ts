import { Language as T } from "@battle-bricks/contracts/common/v1/language_pb";

export function listLanguages(): T[] {
	return Object.values(T).filter((v) => typeof v === "number" && v > 0) as T[];
}

export function getLanguageName(t: T): string {
	switch (t) {
		case T.EN:
			return /* @wc-ignore */ "English";
		case T.DE:
			return /* @wc-ignore */ "German (Deutsch)";
		case T.BG:
			return /* @wc-ignore */ "Bulgarian (Български)";
		default:
			return /* @wc-include */ "Unknown";
	}
}

export function localeToLanguage(locale: string): T {
	switch (locale) {
		case "de":
			return T.DE;
		case "bg":
			return T.BG;
		default:
			return T.EN;
	}
}

export function languageToLocale(language: T): string {
	switch (language) {
		case T.DE:
			return "de";
		case T.BG:
			return "bg";
		default:
			return "en";
	}
}
