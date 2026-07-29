import type { Loadout } from "$lib/loadouts";
import {
	type Ammunition,
	type RangeWeapon,
	type MeleeWeapon,
} from "$lib/items";

interface AttackStat {
	weapon: RangeWeapon | MeleeWeapon;
	ammunition?: Ammunition;
	b1r: number;
	b1o: number;
	b2: number;
	damage: string;
}

export function getAttackStats(attacker: Loadout, defender: Loadout) {
	const attackerWeapons = Array.from(
		new Map(
			[...attacker.unit.items, ...attacker.items]
				.filter(
					(item) =>
						item.type === "RANGE_WEAPON" || item.type === "MELEE_WEAPON",
				)
				.map((item) => [item.id, item]),
		).values(),
	);

	const attackerAmmunitions = Array.from(
		new Map(
			[...attacker.unit.items, ...attacker.items]
				.filter((item) => item.type === "AMMUNITION")
				.map((item) => [item.id, item]),
		).values(),
	);

	return attackerWeapons.reduce<AttackStat[]>((stats, weapon) => {
		if (weapon.type === "RANGE_WEAPON") {
			attackerAmmunitions.forEach((ammunition) => {
				if (ammunition.ammunitionType === weapon.ammunitionType) {
					stats.push({
						weapon: weapon,
						ammunition: ammunition,
						b1r: Math.min(
							Math.max(7 - attacker.unit.marksmanship - defender.unit.size, 0),
							6,
						),
						b1o: Math.min(
							Math.max(8 - attacker.unit.marksmanship - defender.unit.size, 0),
							6,
						),
						b2: Math.min(
							Math.max(
								3 - ammunition.armorPiercing + defender.unit.armorClass,
								0,
							),
							6,
						),
						damage: ammunition.damage,
					});
				}
			});
		} else {
			stats.push({
				weapon: weapon,
				b1r: Math.min(
					Math.max(7 - attacker.unit.meleeAbility - defender.unit.size, 0),
					6,
				),
				b1o: 0,
				b2: Math.min(
					Math.max(3 - weapon.armorPiercing + defender.unit.armorClass, 0),
					6,
				),
				damage: weapon.damage,
			});
		}

		return stats;
	}, []);
}
