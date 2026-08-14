import type {
	Blaster,
	Cannon,
	Launcher,
	MeleeWeapon,
} from "@battle-bricks/contracts/catalogue/v1/item_pb";
import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
import { isGameItem, type GameItem, type GameLoadout } from "$lib/game";
import type { DiceRoll } from "@battle-bricks/contracts/common/v1/dice-roll_pb";

export interface BlasterAction {
	type: "BLASTER";
	name: string;
	weapon: Blaster;
	ammunition: GameItem<"blasterBolt">;
	b1r: number;
	b2?: number;
	damage: DiceRoll;
}

export interface CannonAction {
	type: "CANNON";
	name: string;
	weapon: Cannon;
	ammunition: GameItem<"cannonShell">;
	b1r: number;
	b2?: number;
	damage: DiceRoll;
}

export interface LauncherAction {
	type: "LAUNCHER";
	name: string;
	weapon: Launcher;
	ammunition: GameItem<"launcherRocket">;
	b1r: number;
	b2?: number;
	damage: DiceRoll;
}

export interface MeleeAction {
	type: "MELEE";
	name: string;
	weapon: MeleeWeapon;
	b1r: number;
	b2: number;
	damage: DiceRoll;
}

export interface Action {
	type: "ACTION";
	name: string;
	description?: string;
}

const HitTable: Record<string, number> = {
	"1/1": 8,
	"1/2": 7,
	"1/3": 6,
	"1/4": 6,
	"1/5": 5,

	"2/1": 7,
	"2/2": 6,
	"2/3": 5,
	"2/4": 5,
	"2/5": 4,

	"3/1": 6,
	"3/2": 5,
	"3/3": 4,
	"3/4": 4,
	"3/5": 3,

	"4/1": 5,
	"4/2": 4,
	"4/3": 3,
	"4/4": 3,
	"4/5": 2,

	"5/1": 4,
	"5/2": 3,
	"5/3": 2,
	"5/4": 2,
	"5/5": 1,
};

const PierceTable: Record<string, number> = {
	"1/5": 8,
	"1/4": 8,
	"1/3": 7,
	"1/2": 6,
	"1/1": 5,

	"2/5": 7,
	"2/4": 6,
	"2/3": 6,
	"2/2": 5,
	"2/1": 4,

	"3/5": 6,
	"3/4": 5,
	"3/3": 5,
	"3/2": 4,
	"3/1": 3,

	"4/5": 5,
	"4/4": 4,
	"4/3": 3,
	"4/2": 3,
	"4/1": 2,

	"5/5": 4,
	"5/4": 3,
	"5/3": 2,
	"5/2": 2,
	"5/1": 1,
};

export function getActions(attacker: GameLoadout, defender?: Unit) {
	const actions = [
		...attacker.unit!.items,
		...attacker.items.map((i) => i.item),
	].reduce<
		(BlasterAction | CannonAction | LauncherAction | MeleeAction | Action)[]
	>((stats, item) => {
		if (item.details.case === "blaster") {
			if (attacker.unit.marksmanship && defender) {
				attacker.items.forEach((ammunition) => {
					if (isGameItem(ammunition, "blasterBolt")) {
						stats.push({
							type: "BLASTER",
							name: `${item.name} with ${ammunition.item.name}`,
							weapon: item.details.value as Blaster,
							ammunition: ammunition,
							b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
							b2: PierceTable[
								`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
							],
							damage: ammunition.item.details.value.damage!,
						});
					}
				});
			}
		} else if (item.details.case === "cannon") {
			if (attacker.unit.marksmanship && defender) {
				attacker.items.forEach((ammunition) => {
					if (isGameItem(ammunition, "cannonShell")) {
						stats.push({
							type: "CANNON",
							name: `${item.name} with ${ammunition.item.name}`,
							weapon: item.details.value as Cannon,
							ammunition: ammunition,
							b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
							b2: PierceTable[
								`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
							],
							damage: ammunition.item.details.value.damage!,
						});
					}
				});
			}
		} else if (item.details.case === "launcher") {
			if (attacker.unit.marksmanship && defender) {
				attacker.items.forEach((ammunition) => {
					if (isGameItem(ammunition, "launcherRocket")) {
						stats.push({
							type: "LAUNCHER",
							name: `${item.name} with ${ammunition.item.name}`,
							weapon: item.details.value as Launcher,
							ammunition: ammunition,
							b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
							b2: PierceTable[
								`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
							],
							damage: ammunition.item.details.value.damage!,
						});
					}
				});
			}
		} else if (item.details.case === "meleeWeapon") {
			if (attacker.unit.meleeAbility && defender) {
				stats.push({
					type: "MELEE",
					name: item.name,
					weapon: item.details.value as MeleeWeapon,
					b1r: 6 - attacker.unit.meleeAbility,
					b2: PierceTable[
						`${item.details.value.armorPiercing}/${defender.armorClass}`
					],
					damage: item.details.value.damage!,
				});
			}
		}
		return stats;
	}, []);

	attacker.unit.actions.forEach((action) => {
		actions.push({
			type: "ACTION",
			name: action.name,
			description: action.description,
		});
	});

	return actions;
}
