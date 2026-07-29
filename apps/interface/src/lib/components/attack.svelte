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
			{#if stat.weapon.type === "RANGE_WEAPON"}
				<b>{stat.weapon.name} with {stat.ammunition!.name}</b>
				<div>Fire Rate: {stat.weapon.fireRate}</div>
				<div>Range: {stat.weapon.range.min}-{stat.weapon.range.max}m</div>
				<div>To hit: {stat.b1r} - ({stat.b1o} with cover)</div>
				<div>To pierce armor: {stat.b2}</div>
				<div>Damage: {stat.damage}</div>
			{:else}
				<b>{stat.weapon.name}</b>
				<div>Attack Speed: {stat.weapon.attackSpeed}</div>
				<div>To hit: {stat.b1r}</div>
				<div>To pierce armor: {stat.b2}</div>
				<div>Damage: {stat.damage}</div>
			{/if}
		</li>
	{/each}
</ul>

<style lang="scss">
	ul {
		display: grid;
		box-sizing: border-box;
		width: 100%;
		height: fit-content;
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
