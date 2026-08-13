<script lang="ts">
	import { getFactionName } from "$lib/faction";
	import type { Squad } from "$lib/squad";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import StatRow from "./stat-row.svelte";
	import StatTable from "./stat-table.svelte";
	import BrickCard from "./brick-card.svelte";

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

<BrickCard {onclick}>
	<section>
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
	</section>
</BrickCard>

<style lang="scss">
	@use "sass:color";

	section {
		width: 20rem;
		height: 8.2rem;
		display: flex;
		flex-direction: column;
		cursor: pointer;
	}

	h3 {
		box-sizing: border-box;
		margin: 0;
		font-size: 1.1rem;
		font-weight: 600;
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
