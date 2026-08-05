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

	$inspect(game);
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
			<h2>Your Squad</h2>
			<div>
				<ul>
					{#each game.attacker.loadouts as loadout, i}
						<GameLoadoutItem
							{loadout}
							faction={game.attacker.faction}
							selected={i === selectedAttacker}
							onclick={() => (selectedAttacker = i)}
						/>
					{/each}
				</ul>
			</div>
		</section>
		<section class="defender">
			<h2>Enemy</h2>
			<div>
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
			</div>
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
			"header header header" 1rem
			"squad defender attack" 1fr /
			auto auto 1fr;
		gap: 1.2rem;
		padding: 1rem;
		height: 100dvh;
		width: 100dvw;
	}

	header {
		grid-area: header;
	}

	.squad {
		grid-area: squad;
		display: grid;
		margin: 0;
		padding: 0;
		border-radius: 0.8rem;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		overflow: hidden;

		h2 {
			margin: 0;
			padding: 1rem;
			font-size: 1.4rem;
			font-weight: 600;
		}

		div {
			max-height: 100%;
			max-width: 100%;
			overflow: hidden scroll;
			padding: 0.5rem 1rem 1rem 1rem;
		}
	}

	.defender {
		grid-area: defender;
		display: grid;
		margin: 0;
		padding: 0;
		border-radius: 0.8rem;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		overflow: hidden;

		h2 {
			margin: 0;
			padding: 1rem;
			font-size: 1.4rem;
			font-weight: 600;
		}

		div {
			max-height: 100%;
			max-width: 100%;
			overflow: hidden scroll;
			padding: 0.5rem 1rem 1rem 1rem;
		}
	}

	.attack {
		grid-area: attack;
		margin: 0;
		padding: 0;
	}

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		height: max-content;
		width: max-content;
	}

	section {
		max-height: 100%;
		overflow: scroll;
		padding: 1rem;
	}
</style>
