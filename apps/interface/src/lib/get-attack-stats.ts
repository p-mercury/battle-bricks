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
				ammunitions.forEach((ammunition) => {
					if (ammunition.ammunitionType === weapon.ammunitionType) {
						stats.push({
							type: "RANGE",
							weapon: weapon,
							ammunition: ammunition,
							b1r: Math.min(
								Math.max(8 - attacker.unit.marksmanship - defender.size, 0),
								6,
							),
							b2:
								ammunition.ammunitionType !== "ROCKET"
									? Math.min(
											Math.max(
												3 - ammunition.armorPiercing + defender.armorClass,
												0,
											),
											6,
										)
									: undefined,
							damage: ammunition.damage,
						});
					}
				});
			} else {
				stats.push({
					type: "MELEE",
					weapon: weapon,
					b1r: Math.min(
						Math.max(8 - attacker.unit.meleeAbility - defender.size, 0),
						6,
					),
					b2: Math.min(
						Math.max(3 - weapon.armorPiercing + defender.armorClass, 0),
						6,
					),
					damage: weapon.damage,
				});
			}

			return stats;
		},
		[],
	);
}
