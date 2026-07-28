<script lang="ts">
	import Attack from "$lib/components/attack.svelte";
	import Loadout from "$lib/components/loadout.svelte";
	import { loadouts } from "$lib/loadouts";

	let selectedAttacker = $state<string>();
	let selectedDefender = $state<string>();
</script>

<svelte:head>
	<title>Battle Bricks</title>
</svelte:head>

<div class="wrapper">
	<ul>
		{#each Object.values(loadouts) as loadout}
			<Loadout
				{loadout}
				selected={selectedAttacker === loadout.id}
				onclick={() => {
					if (selectedAttacker !== loadout.id) {
						selectedAttacker = loadout.id;
					} else {
						selectedAttacker = undefined;
					}
				}}
			/>
		{/each}
	</ul>
	<ul>
		{#each Object.values(loadouts) as loadout}
			<Loadout
				{loadout}
				selected={selectedDefender === loadout.id}
				onclick={() => {
					if (selectedDefender !== loadout.id) {
						selectedDefender = loadout.id;
					} else {
						selectedDefender = undefined;
					}
				}}
			/>
		{/each}
	</ul>
	{#if selectedAttacker && selectedDefender}
		<Attack
			attacker={loadouts[selectedAttacker]}
			defender={loadouts[selectedDefender]}
		/>
	{/if}
</div>

<style lang="scss">
	.wrapper {
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
	}

	ul {
		margin: 0;
		padding: 0;
		display: grid;
		grid-auto-rows: max-content;
		width: 10rem;
		gap: 0.5rem;
		height: 100%;
	}
</style>
