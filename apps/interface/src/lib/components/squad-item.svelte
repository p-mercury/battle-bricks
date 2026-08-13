<script lang="ts">
	import { getFactionName } from "$lib/faction";
	import type { Squad } from "$lib/squad";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import StatRow from "./stat-row.svelte";
	import StatTable from "./stat-table.svelte";

	let {
		squad,
		loadouts,
		onclick = () => {},
	}: {
		squad: Squad;
		loadouts: { [k: string]: Loadout };
		onclick?: () => void;
	} = $props();

	let burdget = $derived(
		squad.loadouts.reduce(
			(a, i) =>
				a +
				loadouts[i].unit!.price +
				loadouts[i].items.reduce((total, item) => total + item.price, 0),
			0,
		),
	);
</script>

<li {onclick}>
	<h3>{getFactionName(squad.faction)} {squad.name}</h3>
	<div class="content">
		<StatTable>
			<StatRow>
				Faction:
				<div>{getFactionName(squad.faction)}</div>
			</StatRow>
			<StatRow>
				Price:
				<div>{burdget}/1500c</div>
			</StatRow>
			<StatRow>
				Units:
				<div>{squad.loadouts.length}</div>
			</StatRow>
		</StatTable>
	</div>
</li>

<style lang="scss">
	@use "sass:color";

	li {
		box-sizing: border-box;
		width: 100%;
		height: 8.2rem;
		display: flex;
		flex-direction: column;
		margin: 0;
		padding: 0;
		box-shadow: inset 0 0 0 4px #636669;
		background-color: white;
		border-radius: 0.5rem;
		overflow: hidden;
		cursor: pointer;

		&:hover {
			box-shadow: inset 0 0 0 4px color.adjust(#636669, $lightness: -5%);

			h3 {
				background-color: color.adjust(#636669, $lightness: -5%);
			}
		}

		&.selected {
			box-shadow: inset 0 0 0 4px #2452b6;

			h3 {
				background-color: #2452b6;
			}
		}
	}

	h3 {
		box-sizing: border-box;
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
		background-color: #636669;
		color: white;
		width: 100%;
		padding: 0.4rem 0.5rem 0.2rem 0.5rem;
	}

	.content {
		display: flex;
		flex: 1;
		min-height: 0;
		padding: 0 0.2rem 3px;
		gap: 0.6rem;
	}
</style>
