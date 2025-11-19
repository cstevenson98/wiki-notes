<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import { onMount } from 'svelte';
	import { api } from '$lib/utils/api';
	import { pagesStore } from '$lib/stores/pages';

	let { children } = $props();

	onMount(async () => {
		try {
			const pages = await api.getAllPages();
			pagesStore.set(pages);
		} catch (error) {
			console.error('Failed to load pages:', error);
		}
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Wiki Notes</title>
</svelte:head>

<div class="min-h-screen bg-white">
	<header class="border-b border-gray-200">
		<div class="max-w-3xl mx-auto px-4 py-4">
			<div class="flex items-center justify-between gap-4">
				<a href="/" class="text-xl font-bold hover:text-gray-600 transition-colors">
					Wiki Notes
				</a>
				<SearchBar />
			</div>
		</div>
	</header>

	<main class="max-w-3xl mx-auto px-4 py-8">
		{@render children()}
	</main>
</div>
