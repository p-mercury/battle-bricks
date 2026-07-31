<script lang="ts">
	import Attack from "$lib/components/attack.svelte";
	import UnitItem from "$lib/components/unit-item.svelte";
	import { browser } from "$app/env";
	import type { Game } from "$lib/game";
	import type { Unit } from "$lib/units";
	import GameLoadoutItem from "./game-loadout-item.svelte";

	let game = $state<Game>();
	let initiative = $state(0);
	let selectedDefender = $state<Unit>();
	let selectedAttacker = $derived(
		game?.attacker.loadouts.find((l) => l.initiative === initiative),
	);

	$effect(() => {
		if (browser) {
			let item = localStorage.getItem("GAME");
			if (!item) {
				item = "{}";
			}
			game = JSON.parse(item);
		}
	});

	$effect(() => {
		if (browser) {
			localStorage.setItem("GAME", JSON.stringify(game));
		}
	});
</script>

{#if game}
	<div>
		<button
			disabled={initiative < 1}
			onclick={(event) => {
				event.preventDefault();
				event.stopPropagation();
				initiative--;
			}}>-</button
		>
		{initiative}
		<button
			disabled={initiative > 500}
			onclick={(event) => {
				event.preventDefault();
				event.stopPropagation();
				initiative++;
			}}>+</button
		>
	</div>
	<div class="wrapper">
		<section>
			<h2>Squad</h2>
			<ul>
				{#each game.attacker.loadouts as loadout}
					<GameLoadoutItem
						{loadout}
						bind:initiative={loadout.initiative}
						bind:health={loadout.unit.health}
						selected={loadout.initiative === initiative}
					/>
				{/each}
			</ul>
		</section>
		<section>
			<h2>Enemy</h2>
			<ul>
				{#each game.defender.units as unit}
					<UnitItem
						{unit}
						selected={selectedDefender?.id === unit.id}
						onclick={() => {
							if (selectedDefender?.id !== unit.id) {
								selectedDefender = unit;
							} else {
								selectedDefender = undefined;
							}
						}}
					/>
				{/each}
			</ul>
		</section>
		{#if selectedAttacker && selectedDefender}
			<Attack attacker={selectedAttacker} defender={selectedDefender} />
		{/if}
	</div>
{/if}

<style lang="scss">
	.wrapper {
		box-sizing: border-box;
		display: grid;
		grid-template-columns: auto auto auto 1fr;
		gap: 2rem;
		height: 100dvh;
	}

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
		height: max-content;
		width: max-content;
	}

	section {
		max-height: 100%;
		overflow: scroll;
		padding: 1rem;
	}
</style>
