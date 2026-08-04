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

<div class="color-picker" bind:this={containerEl} style="--color: {value};">
	<button
		type="button"
		class="brick"
		onclick={() => (open = !open)}
		aria-label="Select color"
		aria-expanded={open}
	></button>

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
		background-image:
			linear-gradient(
				to bottom,
				rgb(255 255 255 / 0.18) 0 6%,
				transparent 6% 92%,
				rgb(0 0 0 / 0.22) 92% 100%
			),
			url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='32' height='32'%3E%3Cdefs%3E%3ClinearGradient id='edge' x1='0' y1='0' x2='1' y2='1'%3E%3Cstop offset='0' stop-color='%23fff' stop-opacity='.46'/%3E%3Cstop offset='.46' stop-color='%23fff' stop-opacity='.16'/%3E%3Cstop offset='.54' stop-color='%23000' stop-opacity='.13'/%3E%3Cstop offset='1' stop-color='%23000' stop-opacity='.39'/%3E%3C/linearGradient%3E%3C/defs%3E%3Ccircle cx='16' cy='16' r='9.15' fill='none' stroke='url(%23edge)' stroke-width='1.6'/%3E%3C/svg%3E");
		background-size:
			100% 100%,
			auto 100%;
		background-repeat: no-repeat, repeat-x;
		background-position: center;
		transition: transform 0.1s ease;
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
