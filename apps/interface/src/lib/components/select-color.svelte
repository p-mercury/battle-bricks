<script lang="ts">
	import { colors } from "$lib/color";
	import Brick from "./brick.svelte";

	let { value = $bindable() }: { value: string } = $props();

	let open = $state(false);
	let containerEl: HTMLDivElement;

	function selectColor(hex: string) {
		value = hex;
		open = false;
	}

	function handleClickOutside(event: MouseEvent) {
		if (containerEl && !containerEl.contains(event.target as Node)) {
			open = false;
		}
	}

	$effect(() => {
		if (open) {
			document.addEventListener("click", handleClickOutside);
			return () => document.removeEventListener("click", handleClickOutside);
		}
	});
</script>

<div class="color-picker" bind:this={containerEl} style="--color: {value};">
	<button
		type="button"
		class="brick"
		onclick={() => (open = !open)}
		aria-label="Select color"
		aria-expanded={open}
	>
		<Brick />
	</button>
	{#if open}
		<div class="dropdown">
			{#each colors as color}
				<button
					type="button"
					class="brick option"
					style="--color: {color.hex};"
					onclick={() => selectColor(color.hex)}
					aria-label={color.name}
					title={color.name}
					class:selected={value === color.hex}
				></button>
			{/each}
		</div>
	{/if}
</div>

<style lang="scss">
	.color-picker {
		position: relative;
		display: inline-block;
		height: 100%;
		width: 100%;
	}

	button {
		all: unset;
		cursor: pointer;
	}

	.brick {
		display: block;
		width: 100%;
		height: 100%;
		background-color: var(--color);
	}

	.dropdown {
		position: absolute;
		top: calc(100% + 8px);
		left: 0;
		display: grid;
		grid-template-columns: repeat(4, 1fr);
		gap: 8px;
		padding: 10px;
		background: white;
		border-radius: 8px;
		box-shadow: 0 2px 10px rgba(0, 0, 0, 0.15);
		z-index: 10;
	}

	.option {
		width: 28px;
		height: 28px;
		outline: 2px solid transparent;
		outline-offset: 2px;
	}

	.option:hover {
		transform: scale(1.1);
	}

	.option.selected {
		outline-color: #333;
	}
</style>
