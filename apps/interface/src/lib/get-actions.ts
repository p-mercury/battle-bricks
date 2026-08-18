import type {
	Blaster,
	Cannon,
	Launcher,
	MeleeWeapon,
} from "@battle-bricks/contracts/catalogue/v1/item_pb";
import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
import { isGameItem, type GameItem, type Loadout } from "$lib/game";
import type { DiceRoll } from "@battle-bricks/contracts/common/v1/dice-roll_pb";

export interface BlasterAction {
	type: "BLASTER";
	name: string;
	weapon: Blaster;
	ammunition: GameItem<"blasterBolt">;
	toHit: number;
	toPierce: number;
	damage: DiceRoll;
	damageChance: number;
}

export interface CannonAction {
	type: "CANNON";
	name: string;
	weapon: Cannon;
	ammunition: GameItem<"cannonShell">;
	toHit: number;
	toPierce: number;
	damage: DiceRoll;
	damageChance: number;
}

export interface LauncherAction {
	type: "LAUNCHER";
	name: string;
	weapon: Launcher;
	ammunition: GameItem<"launcherRocket">;
	toHit: number;
	toPierce: number;
	damage: DiceRoll;
	damageChance: number;
}

export interface MeleeAction {
	type: "MELEE";
	name: string;
	weapon: MeleeWeapon;
	toHit: number;
	toPierce: number;
	damage: DiceRoll;
	damageChance: number;
}

export interface Action {
	type: "ACTION";
	name: string;
	description?: string;
}

const HitTable: Record<string, number> = {
	"1/1": 8,
	"1/2": 8,
	"1/3": 7,
	"1/4": 6,
	"1/5": 5,

	"2/1": 7,
	"2/2": 6,
	"2/3": 6,
	"2/4": 5,
	"2/5": 4,

	"3/1": 6,
	"3/2": 5,
	"3/3": 5,
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

const MeleeHitTable: Record<string, number> = {
	"1/5": 8,
	"1/4": 8,
	"1/3": 7,
	"1/2": 6,
	"1/1": 5,
	"1/0": 2,

	"2/5": 7,
	"2/4": 6,
	"2/3": 6,
	"2/2": 5,
	"2/1": 4,
	"2/0": 1,

	"3/5": 6,
	"3/4": 5,
	"3/3": 5,
	"3/2": 4,
	"3/1": 3,
	"3/0": 1,

	"4/5": 5,
	"4/4": 4,
	"4/3": 3,
	"4/2": 3,
	"4/1": 2,
	"4/0": 1,

	"5/5": 4,
	"5/4": 3,
	"5/3": 2,
	"5/2": 2,
	"5/1": 1,
	"5/0": 1,
};

export function getActions(attacker: Loadout, defender?: Unit) {
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
							toHit: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
							toPierce:
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
							damage: ammunition.item.details.value.damage!,
							damageChance: damageChance(
								(item.details.value as Blaster).fireRate,
								HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
								ammunition.item.details.value.damage!.count,
								ammunition.item.details.value.damage!.sides,
								ammunition.item.details.value.damage!.modifier,
							),
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
							toHit: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
							toPierce:
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
							damage: ammunition.item.details.value.damage!,
							damageChance: damageChance(
								(item.details.value as Cannon).fireRate,
								HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
								ammunition.item.details.value.damage!.count,
								ammunition.item.details.value.damage!.sides,
								ammunition.item.details.value.damage!.modifier,
							),
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
							toHit:
								HitTable[
									`${Math.min(attacker.unit.marksmanship! + ammunition.item.details.value.precision, 5)}/${defender.size}`
								],
							toPierce:
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
							damage: ammunition.item.details.value.damage!,
							damageChance: damageChance(
								(item.details.value as Launcher).fireRate,
								HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								PierceTable[
									`${ammunition.item.details.value.armorPiercing}/${defender.armorClass}`
								],
								ammunition.item.details.value.damage!.count,
								ammunition.item.details.value.damage!.sides,
								ammunition.item.details.value.damage!.modifier,
							),
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
					toHit:
						MeleeHitTable[
							`${attacker.unit.meleeAbility}/${defender.meleeAbility || 0}`
						],
					toPierce:
						PierceTable[
							`${(item.details.value as MeleeWeapon).armorPiercing}/${defender.armorClass}`
						],
					damage: item.details.value.damage!,
					damageChance: damageChance(
						(item.details.value as MeleeWeapon).attackSpeed,
						MeleeHitTable[
							`${attacker.unit.meleeAbility}/${defender.meleeAbility || 0}`
						],
						PierceTable[
							`${(item.details.value as MeleeWeapon).armorPiercing}/${defender.armorClass}`
						],
						(item.details.value as MeleeWeapon).damage!.count,
						(item.details.value as MeleeWeapon).damage!.sides,
						(item.details.value as MeleeWeapon).damage!.modifier,
					),
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

function damageChance(
	numberOfDice: number,
	toHit: number,
	toPierce: number,
	numberOfDamageDice: number,
	damageDie: number,
	damageModifier: number,
): number {
	const passChance = (target: number, sides: number): number => {
		const successfulFaces = sides - target + 1;
		return Math.max(0, Math.min(1, successfulFaces / sides));
	};

	const hitChance = passChance(toHit, 8);
	const armourChance = passChance(toPierce, 8);

	// Damage succeeds when roll + modifier >= 1
	const damageTarget = 1 - damageModifier;
	const damageChance = passChance(damageTarget, damageDie);

	// At least one of the generated damage dice succeeds
	const damageAfterPiercing =
		1 - Math.pow(1 - damageChance, numberOfDamageDice);

	// One attack passes all stages
	const damageChancePerAttack = hitChance * armourChance * damageAfterPiercing;

	// At least one of all attacks causes damage
	return (1 - Math.pow(1 - damageChancePerAttack, numberOfDice)) * 100;
}
