<script lang="ts">
	import { getAttackStats } from "$lib/get-attack-stats";
	import type { Loadout } from "$lib/loadouts";
	import type { Unit } from "$lib/units";
	import ActionItem from "./action-item.svelte";

	let { attacker, defender }: { attacker: Loadout; defender: Unit } = $props();

	let actions = $derived(getAttackStats(attacker, defender));
</script>

<ul>
	{#each actions as action}
		<ActionItem {action} />
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
</style>
