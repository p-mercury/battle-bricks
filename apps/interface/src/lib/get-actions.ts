import type { Loadout } from "$lib/loadouts";
import type {
	BoltAmmunition,
	ShellAmmunition,
	RocketAmmunition,
	RangeWeapon,
	MeleeWeapon,
	Action,
} from "$lib/items";
import type { Unit } from "$lib/units";

export interface RangeBoltAction {
	type: "RANGE_BOLT";
	weapon: RangeWeapon;
	ammunition: BoltAmmunition[];
	b1r: number;
	b2?: number;
	damage: string;
}

export interface RangeShellAction {
	type: "RANGE_SHELL";
	weapon: RangeWeapon;
	ammunition: ShellAmmunition[];
	b1r: number;
	b2?: number;
	damage: string;
}

export interface RangeRocketAction {
	type: "RANGE_ROCKET";
	weapon: RangeWeapon;
	ammunition: RocketAmmunition[];
	b1r: number;
	b2?: number;
	damage: string;
}

export interface MeleeAction {
	type: "MELEE";
	weapon: MeleeWeapon;
	b1r: number;
	b2: number;
	damage: string;
}

const HitTable: Record<string, number> = {
	"1/1": 6,
	"1/2": 5,
	"1/3": 5,
	"1/4": 4,
	"2/1": 5,
	"2/2": 4,
	"2/3": 4,
	"2/4": 3,
	"3/1": 4,
	"3/2": 3,
	"3/3": 3,
	"3/4": 2,
	"4/1": 3,
	"4/2": 2,
	"4/3": 2,
	"4/4": 1,
};

const PierceTable: Record<string, number> = {
	"1/1": 3,
	"1/2": 4,
	"1/3": 5,
	"1/4": 6,
	"2/1": 2,
	"2/2": 3,
	"2/3": 4,
	"2/4": 5,
	"3/1": 2,
	"3/2": 2,
	"3/3": 3,
	"3/4": 4,
	"4/1": 1,
	"4/2": 2,
	"4/3": 2,
	"4/4": 3,
};

export function getActions(attacker: Loadout, defender: Unit) {
	const ammunitions = Array.from(
		new Map(
			[...attacker.unit.items, ...attacker.items]
				.filter(
					(item) =>
						item.type === "BOLT_AMMUNITION" ||
						item.type === "SHELL_AMMUNITION" ||
						item.type === "ROCKET_AMMUNITION",
				)
				.map((item) => [item.id, item]),
		).values(),
	);

	return [...attacker.unit.items, ...attacker.items].reduce<
		(
			| RangeBoltAction
			| RangeShellAction
			| RangeRocketAction
			| MeleeAction
			| Action
		)[]
	>((stats, item) => {
		if (item.type === "RANGE_WEAPON") {
			if (attacker.unit.marksmanship) {
				if (item.ammunitionType === "BOLT") {
					ammunitions.forEach((ammunition) => {
						if (ammunition.ammunitionType === "BOLT") {
							stats.push({
								type: "RANGE_BOLT",
								weapon: item,
								ammunition: [...attacker.unit.items, ...attacker.items].filter(
									(a) => a.id === ammunition.id,
								) as any,
								b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								b2: PierceTable[
									`${ammunition.armorPiercing}/${defender.armorClass}`
								],
								damage: ammunition.damage,
							});
						}
					});
				} else if (item.ammunitionType === "SHELL") {
					ammunitions.forEach((ammunition) => {
						if (ammunition.ammunitionType === "SHELL") {
							stats.push({
								type: "RANGE_SHELL",
								weapon: item,
								ammunition: [...attacker.unit.items, ...attacker.items].filter(
									(a) => a.id === ammunition.id,
								) as any,
								b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								b2: PierceTable[
									`${ammunition.armorPiercing}/${defender.armorClass}`
								],
								damage: ammunition.damage,
							});
						}
					});
				} else if (item.ammunitionType === "ROCKET") {
					ammunitions.forEach((ammunition) => {
						if (ammunition.ammunitionType === "ROCKET") {
							stats.push({
								type: "RANGE_ROCKET",
								weapon: item,
								ammunition: [...attacker.unit.items, ...attacker.items].filter(
									(a) => a.id === ammunition.id,
								) as any,
								b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								b2: PierceTable[
									`${ammunition.armorPiercing}/${defender.armorClass}`
								],
								damage: ammunition.damage,
							});
						}
					});
				}
			}
		} else if (item.type === "MELEE_WEAPON") {
			if (attacker.unit.meleeAbility) {
				stats.push({
					type: "MELEE",
					weapon: item,
					b1r: 6 - attacker.unit.meleeAbility,
					b2: PierceTable[`${item.armorPiercing}/${defender.armorClass}`],
					damage: item.damage,
				});
			}
		} else if (item.type === "ACTION") {
			stats.push(item);
		}
		return stats;
	}, []);
}
