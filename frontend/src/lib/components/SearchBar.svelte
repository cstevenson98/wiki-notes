<script lang="ts">
	import { pagesStore } from '$lib/stores/pages';
	import { goto } from '$app/navigation';
	import Input from './ui/Input.svelte';

	let searchQuery = $state('');
	let showResults = $state(false);
	let searchTimeout: ReturnType<typeof setTimeout>;

	const filteredPages = $derived(
		searchQuery.trim()
			? $pagesStore.filter((page) =>
					page.name.toLowerCase().includes(searchQuery.toLowerCase())
				)
			: []
	);

	function handleInput() {
		clearTimeout(searchTimeout);
		searchTimeout = setTimeout(() => {
			showResults = searchQuery.trim().length > 0;
		}, 200);
	}

	function selectPage(pageId: number) {
		goto(`/page/${pageId}`);
		searchQuery = '';
		showResults = false;
	}

	function handleBlur() {
		// Delay to allow click events to fire
		setTimeout(() => {
			showResults = false;
		}, 200);
	}
</script>

<div class="relative w-full max-w-md">
	<Input
		type="text"
		bind:value={searchQuery}
		placeholder="Search pages..."
		oninput={handleInput}
		onfocus={() => (showResults = searchQuery.trim().length > 0)}
		onblur={handleBlur}
		class="w-full"
	/>

	{#if showResults && filteredPages.length > 0}
		<div
			class="absolute top-full mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg max-h-80 overflow-y-auto z-50"
		>
			{#each filteredPages as page (page.id)}
				<button
					class="w-full text-left px-4 py-2 hover:bg-gray-100 border-b border-gray-100 last:border-b-0"
					onclick={() => selectPage(page.id)}
				>
					<div class="font-medium">{page.name}</div>
				</button>
			{/each}
		</div>
	{:else if showResults && searchQuery.trim() && filteredPages.length === 0}
		<div
			class="absolute top-full mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg p-4 z-50"
		>
			<p class="text-sm text-gray-500">No pages found</p>
		</div>
	{/if}
</div>

