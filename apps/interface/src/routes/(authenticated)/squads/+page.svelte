<script lang="ts">
	import { browser } from "$app/environment";
	import { goto } from "$app/navigation";
	import SquadItem from "$lib/components/squad-item.svelte";
	import type { Squad } from "$lib/squad";
	import type { PageData } from "./$types";

	let { data }: { data: PageData } = $props();

	let squads = $derived.by<{ [key: string]: Squad }>(() => {
		if (browser) {
			let item = localStorage.getItem("SQUADS");
			if (!item) {
				item = "{}";
			}
			return JSON.parse(item) as { [key: string]: Squad };
		} else {
			return {};
		}
	});
</script>

<button onclick={() => goto("/squads/new")}>New Squad</button>
<ul>
	{#each Object.values(squads) as squad}
		<SquadItem
			{squad}
			loadouts={data.loadouts}
			onclick={() => goto(`/squads/${squad.id}`)}
		/>
	{/each}
</ul>

<style lang="scss">
	ul {
		all: unset;
		padding: 1rem;
		display: grid;
		gap: 1rem;
	}
</style>
