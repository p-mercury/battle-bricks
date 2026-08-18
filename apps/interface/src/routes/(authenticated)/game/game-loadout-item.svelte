<script lang="ts">
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import NumberInput from "$lib/components/number-input.svelte";
	import type { Loadout } from "$lib/game";
	import Switcher from "$lib/components/switcher.svelte";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import BrickCard from "$lib/components/brick-card.svelte";

	let {
		faction,
		loadout,
		selected = false,
		onclick = () => {},
	}: {
		faction: Faction;
		loadout: Loadout;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<BrickCard {selected} {onclick} bind:color={loadout.color}>
	<section>
		{#if loadout.image}
			<img
				alt={loadout.name}
				src={loadout.image}
				class:republic={faction === Faction.GALACTIC_REPUBLIC}
				class:separatists={faction === Faction.SEPARATIST_ALLIANCE}
			/>
		{/if}
		<div class="info">
			<h3>{loadout.name}</h3>
		</div>
		<div class="stats">
			<StatTable>
				<StatRow>
					Turn:
					<Switcher
						bind:value={loadout.turnComplete}
						onclick={(event) => event.stopPropagation()}
					/>
				</StatRow>
				<StatRow>
					Health:
					<NumberInput bind:value={loadout.unit.health} />
				</StatRow>
				{#if "inCover" in loadout}
					<StatRow>
						In Cover:
						<Switcher
							bind:value={loadout.inCover}
							onclick={(event) => event.stopPropagation()}
						/>
					</StatRow>
				{/if}
				<StatRow>
					Size:
					<StatBar value={loadout.unit.size} size={5} />
				</StatRow>
				<StatRow>
					Speed:
					<StatBar value={loadout.unit.speed} size={5} />
				</StatRow>
				<StatRow>
					Armor Class:
					<StatBar value={loadout.unit.armorClass} size={5} />
				</StatRow>
				{#if loadout.unit.marksmanship}
					<StatRow>
						Marksmanship:
						<StatBar value={loadout.unit.marksmanship} size={5} red />
					</StatRow>
				{/if}
				{#if loadout.unit.meleeAbility}
					<StatRow>
						Melee Ability:
						<StatBar value={loadout.unit.meleeAbility} size={5} red />
					</StatRow>
				{/if}
			</StatTable>
		</div>
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		padding: 0.6rem;
		display: grid;
		width: 21.55rem;
		grid-template:
			"image stats" 7rem
			"info stats" auto /
			7rem auto;
		gap: 0.6rem;
	}

	img {
		width: 100%;
		height: 100%;
		grid-area: image;
		object-fit: contain;
		border-radius: 0.8rem;
		padding: 0.4rem;
		background: linear-gradient(135deg, #1c1c1c, #4f424f);
		box-sizing: border-box;

		&.republic {
			background: linear-gradient(135deg, #1c1c1c, #4f424f);
		}

		&.separatists {
			background: linear-gradient(135deg, #2b1515, #2d3f54);
		}
	}

	.stats {
		grid-area: stats;
		padding-left: 0.3rem;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0;
			font-size: 1.1rem;
			font-weight: 600;
		}
	}
</style>
