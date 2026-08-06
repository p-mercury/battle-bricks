<script lang="ts">
	import StatBar from "$lib/components/stat-bar.svelte";
	import StatRow from "$lib/components/stat-row.svelte";
	import StatTable from "$lib/components/stat-table.svelte";
	import type { Faction } from "$lib/faction";
	import type { Unit } from "$lib/units";
	import Brick from "./brick.svelte";

	let {
		faction,
		unit,
		selected = false,
		onclick = () => {},
	}: {
		faction?: Faction;
		unit: Unit;
		selected?: Boolean;
		onclick?: () => void;
	} = $props();
</script>

<section {onclick} class:selected>
	<header>
		<Brick />
	</header>
	{#if unit.image}
		<img
			alt={unit.name}
			src={unit.image}
			class:republic={faction === "GALACTIC_REPUBLIC"}
			class:separatists={faction === "SEPARATIST_ALLIANCE"}
		/>
	{/if}
	<div class="info">
		<h3>{unit.name}</h3>
		<StatTable>
			<StatRow>
				Health:
				<div>{unit.health}hp</div>
			</StatRow>
			<StatRow>
				Speed:
				<StatBar value={unit.speed} size={4} />
			</StatRow>
			<StatRow>
				Size:
				<StatBar value={unit.size} size={4} />
			</StatRow>
		</StatTable>
	</div>
	<div class="stats">
		<StatTable>
			<StatRow>
				Armor Class:
				<StatBar value={unit.armorClass} size={4} />
			</StatRow>
			{#if unit.marksmanship}
				<StatRow>
					Marksmanship:
					<StatBar value={unit.marksmanship} size={4} red />
				</StatRow>
			{/if}
			{#if unit.meleeAbility}
				<StatRow>
					Melee Ability:
					<StatBar value={unit.meleeAbility} size={4} red />
				</StatRow>
			{/if}
		</StatTable>
	</div>
</section>

<style lang="scss">
	@use "sass:color";

	section {
		background: white;
		border-radius: 0.8rem;
		padding: 0.5rem;
		overflow: hidden;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		cursor: pointer;
		&.selected {
			outline: 4px solid #2388ff;
			outline-offset: 2px;
		}
		display: grid;
		width: 21.8rem;
		grid-template:
			"color color" 1.25rem
			"image info" 8rem
			"stats stats" 1fr /
			8rem auto;
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

	header {
		grid-area: color;
		margin: -0.5rem -0.5rem 0 -0.5rem;
		background-color: #595d60;
	}

	.stats {
		grid-area: stats;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0 0 0 0.2rem;
			font-size: 1.2rem;
			font-weight: 600;
		}
	}
</style>
