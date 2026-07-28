<script lang="ts">
	import { getAttackStats } from "$lib/get-attack-stats";
	import type { Loadout } from "$lib/loadouts";

	let { attacker, defender }: { attacker: Loadout; defender: Loadout } =
		$props();

	let stats = $derived(getAttackStats(attacker, defender));

	$inspect(stats);
</script>

<ul>
	{#each stats as stat}
		<li>
			<b>{stat.weapon.name} with {stat.ammunition.name}</b>
			<div>Fire Rate: {stat.weapon.fireRate}</div>
			<div>Range: {stat.weapon.range.min}-{stat.weapon.range.max}m</div>
			<div>B1R: {stat.b1r}</div>
			<div>B1O: {stat.b1o}</div>
			<div>B2: {stat.b2}</div>
			<div>Damage: {stat.damage}</div>
		</li>
	{/each}
</ul>

<style lang="scss">
	ul {
		display: grid;
		box-sizing: border-box;
		width: 100%;
		height: 100%;
		margin: 0;
		padding: 0.5rem;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightgray;
		gap: 0.5rem;
	}

	li {
		padding: 0.5rem;
		list-style: none;
		border: 2px solid black;
		border-radius: 0.4rem;
		background-color: lightsalmon;
	}
</style>
