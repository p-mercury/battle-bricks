<script lang="ts">
	import UnitItem from "$lib/components/unit-item.svelte";
	import type { Unit } from "@battle-bricks/contracts/catalogue/v1/unit_pb";
	import { dragHandleZone, type DndEvent } from "svelte-dnd-action";

	let {
		items = $bindable(),
		flipDurationMs,
		handleConsider,
		handleFinalize,
	}: {
		items: { id: string; unit: Unit }[];
		flipDurationMs: number;
		handleConsider: (
			event: CustomEvent<DndEvent<(typeof items)[number]>>,
		) => void;
		handleFinalize: (
			event: CustomEvent<DndEvent<(typeof items)[number]>>,
		) => void;
	} = $props();
</script>

<ul
	use:dragHandleZone={{
		items,
		flipDurationMs,
		morphDisabled: true,
		useCursorForDetection: true,
		transformDraggedElement: (element) => {
			if (!element) return;

			element.style.opacity = "1";
			element.style.background = "#ffffff";
			element.style.borderRadius = "0.6rem";
		},
		dropTargetStyle: {
			outline: "2px dashed #3b82f6",
		},
	}}
	onconsider={handleConsider}
	onfinalize={handleFinalize}
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
