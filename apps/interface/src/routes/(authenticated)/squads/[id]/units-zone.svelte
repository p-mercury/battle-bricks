<script lang="ts">
	import UnitItem from "$lib/components/unit-item.svelte";
	import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
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
		items: { id: string; unit: Unit }[];
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
		<UnitItem unit={item.unit} drag />
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
