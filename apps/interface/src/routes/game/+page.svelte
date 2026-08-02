<script lang="ts">
	import Attack from "$lib/components/attack.svelte";
	import UnitItem from "$lib/components/unit-item.svelte";
	import { browser } from "$app/env";
	import type { Game } from "$lib/game";
	import type { Unit } from "$lib/units";
	import GameLoadoutItem from "./game-loadout-item.svelte";

	let game = $state<Game>();
	let selectedDefender = $state<Unit>();
	let selectedAttacker = $state<number>();

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
	<div class="wrapper">
		<header>
			<button
				onclick={() => {
					game?.attacker.loadouts.forEach((_, i) => {
						game!.attacker.loadouts[i].turnComplete = false;
					});
				}}>Clear Turn Complete</button
			>
		</header>
		<section class="squad">
			<ul>
				{#each game.attacker.loadouts as loadout, i}
					<GameLoadoutItem
						{loadout}
						selected={i === selectedAttacker}
						onclick={() => (selectedAttacker = i)}
					/>
				{/each}
			</ul>
		</section>
		<section class="defender">
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
		<section class="attack">
			{#if selectedAttacker != null && selectedDefender}
				<Attack
					attacker={game.attacker.loadouts[selectedAttacker]}
					defender={selectedDefender}
				/>
			{/if}
		</section>
	</div>
{/if}

<style lang="scss">
	.wrapper {
		box-sizing: border-box;
		display: grid;
		grid-template:
			"header header header" 3rem
			"squad defender attack" 1fr /
			auto auto 1fr;
		gap: 1rem;
		height: 100dvh;
	}

	header {
		grid-area: header;
	}

	.squad {
		grid-area: squad;
	}

	.defender {
		grid-area: defender;
	}

	.attack {
		grid-area: attack;
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
