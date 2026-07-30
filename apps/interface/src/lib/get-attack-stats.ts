import type { Loadout } from "$lib/loadouts";
import {
	type PlasmaAmmunition,
	type RocketAmmunition,
	type SlugAmmunition,
	type RangeWeapon,
	type MeleeWeapon,
} from "$lib/items";
import type { Unit } from "$lib/units";

interface RangeAttackStat {
	type: "RANGE";
	weapon: RangeWeapon;
	ammunition: PlasmaAmmunition | RocketAmmunition | SlugAmmunition;
	b1r: number;
	b2?: number;
	damage: string;
}

interface MeleeAttackStat {
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

export function getAttackStats(attacker: Loadout, defender: Unit) {
	const weapons = Array.from(
		new Map(
			[...attacker.unit.items, ...attacker.items]
				.filter(
					(item) =>
						item.type === "RANGE_WEAPON" || item.type === "MELEE_WEAPON",
				)
				.map((item) => [item.id, item]),
		).values(),
	);

	const ammunitions = Array.from(
		new Map(
			[...attacker.unit.items, ...attacker.items]
				.filter(
					(item) =>
						item.type === "PLASMA_AMMUNITION" ||
						item.type === "ROCKET_AMMUNITION" ||
						item.type === "SLUG_AMMUNITION",
				)
				.map((item) => [item.id, item]),
		).values(),
	);

	return weapons.reduce<(RangeAttackStat | MeleeAttackStat)[]>(
		(stats, weapon) => {
			if (weapon.type === "RANGE_WEAPON") {
				if (attacker.unit.marksmanship) {
					ammunitions.forEach((ammunition) => {
						if (ammunition.ammunitionType === weapon.ammunitionType) {
							stats.push({
								type: "RANGE",
								weapon: weapon,
								ammunition: ammunition,
								b1r: HitTable[`${attacker.unit.marksmanship}/${defender.size}`],
								b2:
									ammunition.ammunitionType !== "ROCKET"
										? PierceTable[
												`${ammunition.armorPiercing}/${defender.armorClass}`
											]
										: undefined,
								damage: ammunition.damage,
							});
						}
					});
				}
			} else {
				if (attacker.unit.meleeAbility) {
					stats.push({
						type: "MELEE",
						weapon: weapon,
						b1r: HitTable[`${attacker.unit.meleeAbility}/${defender.size}`],
						b2: PierceTable[`${weapon.armorPiercing}/${defender.armorClass}`],
						damage: weapon.damage,
					});
				}
			}

			return stats;
		},
		[],
	);
}
