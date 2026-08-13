<script lang="ts">
	import LoadoutItem from "$lib/components/loadout-item.svelte";
	import type { Loadout } from "@battle-bricks/contracts/catalogue/v1/loadout_pb";
	import {
		dragHandleZone,
		TRIGGERS,
		SHADOW_ITEM_MARKER_PROPERTY_NAME,
		type DndEvent,
	} from "svelte-dnd-action";

	let {
		items = $bindable(),
		flipDurationMs,
	}: {
		items: { id: string; loadout: Loadout }[];
		flipDurationMs: number;
	} = $props();

	function handleDndConsider(event: CustomEvent<DndEvent<any>>) {
		const { trigger, id } = event.detail.info;
		if (trigger === TRIGGERS.DRAG_STARTED) {
			const idx = items.findIndex((item) => item.id === id);
			event.detail.items = event.detail.items.filter(
				(item: any) => !item[SHADOW_ITEM_MARKER_PROPERTY_NAME],
			);
			event.detail.items.splice(idx, 0, {
				...$state.snapshot(items[idx]),
				id: crypto.randomUUID(),
			});
			items = event.detail.items;
		} else {
			items = [...items];
		}
	}

	function handleDndFinalize() {
		items = [...items];
	}
</script>

<ul
	use:dragHandleZone={{
		items,
		useCursorForDetection: true,
		flipDurationMs,
		dropTargetStyle: {
			outline: "0px",
		},
	}}
	onconsider={handleDndConsider}
	onfinalize={handleDndFinalize}
>
	{#each items as item (item.id)}
		<LoadoutItem loadout={item.loadout} drag />
	{/each}
</ul>

<style lang="scss">
	ul {
		all: unset;
		display: flex;
		flex-direction: column;
		gap: 1rem;
		padding: 1rem;
		min-height: 100%;
	}
</style>
