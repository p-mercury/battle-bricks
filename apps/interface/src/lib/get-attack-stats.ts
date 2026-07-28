import type { Loadout } from "$lib/loadouts";
import { type Ammunition, type Weapon } from "$lib/items";
import { getLoadoutStats } from "./get-loadout-stats";

interface AttackStat {
	weapon: Weapon;
	ammunition: Ammunition;
	b1r: number;
	b1o: number;
	b2: number;
	damage: string;
}

export function getAttackStats(attacker: Loadout, defender: Loadout) {
	const _attacker = getLoadoutStats(attacker);
	const attackerWeapons = Array.from(
		new Map(
			_attacker.items
				.filter((item) => item.type === "WEAPON")
				.map((item) => [item.id, item]),
		).values(),
	);
	const attackerAmmunitions = Array.from(
		new Map(
			_attacker.items
				.filter((item) => item.type === "AMMUNITION")
				.map((item) => [item.id, item]),
		).values(),
	);

	const _defender = getLoadoutStats(defender);

	return attackerWeapons.reduce<AttackStat[]>((stats, weapon) => {
		stats.push(
			...attackerAmmunitions.map((ammunition) => ({
				weapon: weapon,
				ammunition: ammunition,
				b1r: Math.min(
					Math.max(7 - _attacker.unit.accuracy - _defender.unit.size, 0),
					6,
				),
				b1o: Math.min(
					Math.max(9 - _attacker.unit.accuracy - _defender.unit.size, 0),
					6,
				),
				b2: Math.min(
					Math.max(3 - ammunition.armorPiercing + _defender.unit.armorClass, 0),
					6,
				),
				damage: ammunition.damage,
			})),
		);

		return stats;
	}, []);
}
