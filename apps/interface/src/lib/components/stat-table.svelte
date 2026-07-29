<script lang="ts">
	import type { Snippet } from "svelte";

	let { children }: { children: () => ReturnType<Snippet> } = $props();
</script>

<ul>
	{@render children()}
</ul>

<style lang="scss">
	ul {
		margin: 0;
		display: grid;
		grid-template-columns: auto auto;
		gap: 0.5rem;
		padding: 0.5rem;
		flex-shrink: 0;
	}
</style>
