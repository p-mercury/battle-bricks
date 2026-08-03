<script lang="ts">
	import { colors } from "$lib/color";

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

<div class="color-picker" bind:this={containerEl}>
	<button
		type="button"
		class="color-swatch"
		style="background-color: {value};"
		onclick={() => (open = !open)}
		aria-label="Select color"
		aria-expanded={open}
	></button>
	{#if open}
		<div class="dropdown">
			{#each colors as color}
				<button
					type="button"
					class="color-option"
					style="background-color: {color.hex};"
					onclick={() => selectColor(color.hex)}
					aria-label={color.name}
					title={color.name}
					class:selected={value === color.hex}
				></button>
			{/each}
		</div>
	{/if}
</div>

<style lnag="scss">
	.color-picker {
		position: relative;
		display: inline-block;
		height: 1.8rem;
	}

	button {
		all: unset;
		width: 1.8rem;
		height: 1.8rem;
		border-radius: 50%;
		cursor: pointer;
		padding: 0;
		margin: 0;
	}

	button:hover {
		transform: scale(1.05);
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

	.color-option {
		width: 28px;
		height: 28px;
		border-radius: 50%;
		border: 2px solid #ddd;
		cursor: pointer;
		padding: 0;
		transition: transform 0.1s ease;
	}

	.color-option:hover {
		transform: scale(1.1);
	}

	.color-option.selected {
		border: 2px solid #333;
		box-shadow: 0 0 0 2px white inset;
	}
</style>
