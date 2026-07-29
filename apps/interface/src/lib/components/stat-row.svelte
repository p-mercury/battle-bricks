<script lang="ts">
	import type { Snippet } from "svelte";

	let { children }: { children: () => ReturnType<Snippet> } = $props();
</script>

<li>
	{@render children()}
</li>

<style lang="scss">
	li {
		display: grid;
		list-style: none;
		grid-column: 1 / -1;
		grid-template-columns: subgrid;
	}
</style>
