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
		grid-template-columns: 1fr auto;
		gap: 0.5rem;
		padding: 0;
		flex-shrink: 0;
		height: fit-content;
	}
</style>
