import type { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
import type { Item } from "@battle-bricks/contracts/catalogue/v1/item_pb";
import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";

type ItemCase = Exclude<Item["details"]["case"], undefined>;

type ItemOf<C extends ItemCase> = Item & {
	details: Extract<Item["details"], { case: C }>;
};

export interface GameItem<C extends ItemCase | never = never> {
	quantity: number;
	maxQuantity: number;
	item: [C] extends [never] ? Item : ItemOf<C>;
}

export function isGameItem<C extends ItemCase>(
	gameItem: GameItem,
	itemCase: C,
): gameItem is GameItem & { item: ItemOf<C> } {
	return gameItem.item.details.case === itemCase;
}

export interface GameLoadout {
	id: string;
	image?: string;
	name: string;
	color: string;
	turnComplete: boolean;
	inCover?: boolean;
	unit: Unit;
	items: GameItem[];
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
