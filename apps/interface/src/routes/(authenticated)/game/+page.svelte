<!-- <script lang="ts">
	import UnitItem from "$lib/components/unit-item.svelte";
	import { browser } from "$app/env";
	import type { Game } from "$lib/game";
	import type { Unit } from "$lib/units";
	import GameLoadoutItem from "./game-loadout-item.svelte";
	import ActionItem from "$lib/components/action-item.svelte";
	import { getActions } from "$lib/get-actions";
	import { flip } from "svelte/animate";
	import { quintInOut, quintOut } from "svelte/easing";
	import { fade } from "svelte/transition";

	let game = $state<Game>();
	let selectedDefender = $state<Unit>();
	let selectedAttacker = $state<string>();
	let actions = $derived.by(() => {
		if (game && selectedAttacker) {
			return getActions(
				game.attacker.loadouts[selectedAttacker],
				selectedDefender,
			);
		}
		return [];
	});

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
				onclick={() =>
					Object.values(game!.attacker.loadouts).forEach(({ id }) => {
						if (game!.attacker.loadouts[id].unit.health) {
							game!.attacker.loadouts[id].turnComplete = false;
						}
					})}
			>
				Clear Turn Complete
			</button>
		</header>
		<section class="squad">
			<h2>Your Squad</h2>
			<div>
				<ul>
					{#each Object.values(game!.attacker.loadouts).toSorted((a, b) => Number(a.turnComplete) - Number(b.turnComplete)) as loadout (loadout.id)}
						<li animate:flip={{ duration: 800, easing: quintOut }}>
							<GameLoadoutItem
								{loadout}
								faction={game.attacker.faction}
								selected={loadout.id === selectedAttacker}
								onclick={() => {
									if (selectedAttacker !== loadout.id) {
										selectedAttacker = loadout.id;
									} else {
										selectedAttacker = undefined;
									}
								}}
							/>
						</li>
					{/each}
				</ul>
			</div>
		</section>
		<section class="defender">
			<h2>Enemy</h2>
			<div>
				<ul>
					{#each game.defender.units as unit (unit.id)}
						<li animate:flip={{ duration: 800, easing: quintOut }}>
							<UnitItem
								{unit}
								faction={game.defender.faction}
								selected={selectedDefender?.id === unit.id}
								onclick={() => {
									if (selectedDefender?.id !== unit.id) {
										selectedDefender = unit;
									} else {
										selectedDefender = undefined;
									}
								}}
							/>
						</li>
					{/each}
				</ul>
			</div>
		</section>
		<section class="actions">
			<h2>Actions</h2>
			<div>
				<ul>
					{#each actions as action}
						<li in:fade={{ duration: 600, easing: quintInOut }}>
							<ActionItem {action} />
						</li>
					{/each}
				</ul>
			</div>
		</section>
	</div>
{/if}

<style lang="scss">
	.wrapper {
		box-sizing: border-box;
		display: grid;
		grid-template:
			"header header header header" 1rem
			"squad defender actions ." 1fr /
			auto auto minmax(24rem, auto) 1fr;
		gap: 1rem;
		padding: 0.5rem;
		height: 100dvh;
		width: 100dvw;
	}

	header {
		grid-area: header;
	}

	.squad {
		grid-area: squad;
		display: flex;
		flex-direction: column;
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
			padding: 0.5rem 1rem 0 1rem;
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
		display: flex;
		flex-direction: column;
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
			padding: 0.5rem 1rem 0 1rem;
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

	.actions {
		grid-area: actions;
		display: flex;
		flex-direction: column;
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
			padding: 0.5rem 1rem 0 1rem;
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

	ul {
		margin: 0;
		padding: 0;
		display: flex;
		flex-direction: column;
		gap: 0.8rem;
		height: max-content;
		width: 100%;
	}

	section {
		max-height: 100%;
		overflow: scroll;
		padding: 1rem;
	}

	li {
		all: unset;
	}
</style> -->
