<script>
	import { pagesStore } from '$lib/stores/pages';
</script>

<div class="text-center py-12">
	<h1 class="text-4xl font-bold mb-4">Welcome to Wiki Notes</h1>
	<p class="text-gray-600 mb-8">Your personal wiki-style note-taking application</p>

	{#if $pagesStore.length === 0}
		<p class="text-gray-500">Loading pages...</p>
	{:else}
		<p class="text-gray-500">Redirecting to Home page...</p>
	{/if}
</div>
