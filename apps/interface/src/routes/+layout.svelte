<script lang="ts">
	import type { Snippet } from "svelte";

	let { children }: { children: () => ReturnType<Snippet> } = $props();
</script>

<svelte:head>
	<title>Battle Bricks</title>
</svelte:head>

{@render children()}

<style lang="scss">
	@font-face {
		font-family: "Atkinson";
		src: url(/fonts/atkinson-hyperlegible-next.ttf);
		font-weight: 100 1000;
	}

	:global(body) {
		font-family: "Atkinson";
		margin: 0;
		padding: 0;
	}
</style>
