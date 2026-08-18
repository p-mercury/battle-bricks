<script lang="ts">
	import { navigating } from "$app/state";
	import Loader from "$lib/components/loader.svelte";
	import type { Snippet } from "svelte";

	let { children }: { children: () => ReturnType<Snippet> } = $props();

	let deleayedNavigationTimeout: NodeJS.Timeout | undefined;
	let deleayedNavigation = $state(false);
	$effect(() => {
		if (navigating.type != null) {
			if (!deleayedNavigationTimeout) {
				deleayedNavigationTimeout = setTimeout(() => {
					deleayedNavigation = true;
					clearTimeout(deleayedNavigationTimeout);
					deleayedNavigationTimeout = undefined;
				}, 400);
			}
		} else {
			deleayedNavigation = false;
			clearTimeout(deleayedNavigationTimeout);
			deleayedNavigationTimeout = undefined;
		}
	});
</script>

{#if navigating.type != null}
	{#if deleayedNavigation}
		<Loader />
	{/if}
{:else}
	{@render children()}
{/if}
