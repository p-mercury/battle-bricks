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

interface RangeAttackStat {
	type: "RANGE";
	weapon: RangeWeapon;
	ammunition: BoltAmmunition | ShellAmmunition | RocketAmmunition;
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
		(RangeAttackStat | MeleeAttackStat | Action)[]
	>((stats, item) => {
		console.log(item);
		if (item.type === "RANGE_WEAPON") {
			if (attacker.unit.marksmanship) {
				ammunitions.forEach((ammunition) => {
					if (ammunition.ammunitionType === item.ammunitionType) {
						stats.push({
							type: "RANGE",
							weapon: item,
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
