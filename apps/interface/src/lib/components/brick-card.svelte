<script lang="ts">
	import type { Snippet } from "svelte";
	import Brick from "$lib/components/brick.svelte";
	import SelectColor from "$lib/components/select-color.svelte";
	import { dragHandle } from "svelte-dnd-action";

	let {
		selected,
		color = $bindable(),
		drag = false,
		onclick,
		children,
	}: {
		selected?: Boolean;
		color?: string;
		drag?: Boolean;
		onclick?: () => void;
		children: () => ReturnType<Snippet>;
	} = $props();
</script>

<div {onclick} class="wrapper" class:selected>
	{#if color}
		<div class="brick">
			<SelectColor bind:value={color} />
		</div>
	{:else}
		{#if drag}
			<div class="brick" use:dragHandle>
				<Brick />
			</div>
		{:else}
			<div class="brick">
				<Brick />
			</div>
		{/if}
	{/if}
	<div>{@render children()}</div>
</div>

<style lang="scss">
	.wrapper {
		display: grid;
		background: white;
		border-radius: 0.8rem;
		overflow: hidden;
		box-shadow:
			0 1px 2px rgba(0, 0, 0, 0.1),
			0 4px 8px rgba(0, 0, 0, 0.14),
			0 8px 16px rgba(0, 0, 0, 0.12);
		&.selected {
			outline: 4px solid #2388ff;
			outline-offset: 2px;
		}
		grid-template:
			"brick" 1.75rem
			"icontent" auto /
			min-content;
		width: min-content;
	}

	.brick {
		background-color: #595d60;
	}
</style>
