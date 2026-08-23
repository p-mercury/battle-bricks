<script lang="ts">
	import type { Loadout } from "$lib/game";
	import { Faction } from "@battle-bricks/contracts/catalogue/v1/faction_pb";
	import BrickCard from "$lib/components/brick-card.svelte";
	import Stat from "$lib/components/stat.svelte";
	import NumberInput from "$lib/components/number-input.svelte";
	import Switcher from "$lib/components/switcher.svelte";

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
			<div class="image-container" class:dead={loadout.health === 0}>
				<img
					alt={loadout.name}
					src={loadout.image}
					class:republic={faction === Faction.GALACTIC_REPUBLIC}
					class:separatists={faction === Faction.SEPARATIST_ALLIANCE}
				/>
				{#if loadout.health === 0}
					<span class="dead-label">DEAD</span>
				{/if}
			</div>
		{/if}
		<div class="info">
			<h3>{loadout.name}</h3>
		</div>
		<div class="turn">
			<Switcher bind:value={loadout.turnComplete} />
		</div>
		<div class="stats">
			<Stat label="HP" color="GREEN" value={loadout.unit.health} />
			<Stat label="SZ" color="BLUE" value={loadout.unit.size} />
			<Stat label="GS" color="RED" value={loadout.unit.marksmanship || 0} />
			<Stat label="AC" color="GREEN" value={loadout.unit.armorClass} />
			<Stat label="SP" color="BLUE" value={loadout.unit.speed} />
			<Stat label="MS" color="RED" value={loadout.unit.meleeAbility || 0} />
		</div>
		<div class="trackers">
			<NumberInput
				min={0}
				max={loadout.unit.health}
				bind:value={loadout.health}
			/>
			{#if loadout.inCover != null}
				<Switcher bind:value={loadout.inCover} />
			{/if}
		</div>
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		padding: 0.6rem;
		display: grid;
		width: 18.6rem;
		grid-template:
			"info info turn" 1.3rem
			"image stats stats" 7.5rem
			"trackers trackers trackers" auto /
			7.5rem auto min-content;
		gap: 0.6rem;
	}

	.image-container {
		grid-area: image;
		position: relative;
		width: 100%;
		height: 100%;
		border-radius: 0.8rem;
		overflow: hidden;
		clip-path: inset(0 round 0.8rem);
		isolation: isolate;

		img {
			display: block;
			width: 100%;
			height: 100%;
			object-fit: contain;
			padding: 0.4rem;
			border-radius: inherit;
			background: linear-gradient(135deg, #1c1c1c, #4f424f);
			box-sizing: border-box;
			transition:
				filter 0.2s ease,
				opacity 0.2s ease;

			&.republic {
				background: linear-gradient(135deg, #1c1c1c, #4f424f);
			}

			&.separatists {
				background: linear-gradient(135deg, #2b1515, #2d3f54);
			}
		}

		&.dead img {
			filter: grayscale(1) blur(3px);
			opacity: 0.55;
			transform: scale(1.03);
		}
	}

	.dead-label {
		position: absolute;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #ffffff;
		font-size: 1.5rem;
		font-weight: 900;
		letter-spacing: 0.15rem;
		text-shadow:
			0 2px 4px rgb(0 0 0 / 90%),
			0 0 8px rgb(255 0 0 / 80%);
		pointer-events: none;
	}

	.stats {
		grid-area: stats;
		display: grid;
		grid-template-columns: 1fr 1fr 1fr;
		grid-auto-rows: 1fr;
		gap: 0.4rem;
	}

	.trackers {
		grid-area: trackers;
		display: flex;
		gap: 0.5rem;
	}

	.info {
		grid-area: info;

		h3 {
			margin: 0;
			padding: 0;
			font-size: 1.2rem;
			line-height: 1.2rem;
			font-weight: 600;
		}
	}

	.turn {
		grid-area: turn;
	}
</style>
