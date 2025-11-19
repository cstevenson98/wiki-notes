<script lang="ts">
	import { onMount, tick } from 'svelte';
	import Editor from '$lib/components/Editor.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Dialog from '$lib/components/ui/Dialog.svelte';
	import { api, type Page } from '$lib/utils/api';
	import { pagesStore, currentPageStore } from '$lib/stores/pages';
	import { goto } from '$app/navigation';
	import 'katex/dist/katex.min.css';

	let { data } = $props();

	let page = $state<Page>(data.page);
	let editedContent = $state(data.page.content);
	let isEditing = $state(false);
	let isSaving = $state(false);
	let saveMessage = $state('');
	let showCreateDialog = $state(false);
	let newPageName = $state('');
	let creating = $state(false);
	let contentElement: HTMLDivElement;

	const displayContent = $derived(processWikiLinksForDisplay(page.content));

	$effect(() => {
		currentPageStore.set(page);
	});

	// Render LaTeX when content changes, page changes, or switching from edit to view
	$effect(() => {
		// Depend on displayContent and isEditing to trigger on any change
		const content = displayContent;
		const editing = isEditing;
		
		if (!editing && contentElement && content) {
			tick().then(() => {
				renderMathInElement();
			});
		}
	});

	async function renderMathInElement() {
		if (typeof window !== 'undefined' && contentElement) {
			const katex = await import('katex');
			const renderMathInElement = (await import('katex/dist/contrib/auto-render.mjs')).default;
			
			renderMathInElement(contentElement, {
				delimiters: [
					{ left: '$$', right: '$$', display: true },
					{ left: '$', right: '$', display: false },
					{ left: '\\(', right: '\\)', display: false },
					{ left: '\\[', right: '\\]', display: true }
				],
				throwOnError: false
			});
		}
	}

	async function saveChanges() {
		if (!page) return;

		isSaving = true;
		saveMessage = '';
		try {
			const updatedPage = await api.updatePage(page.id, {
				content: editedContent
			});

			page = updatedPage;
			isEditing = false;
			saveMessage = 'Saved successfully!';

			// Update pages store
			pagesStore.update((pages) =>
				pages.map((p) => (p.id === updatedPage.id ? updatedPage : p))
			);

			setTimeout(() => {
				saveMessage = '';
			}, 3000);
		} catch (error) {
			saveMessage = 'Failed to save: ' + (error as Error).message;
		} finally {
			isSaving = false;
		}
	}

	function cancelEdit() {
		editedContent = page.content;
		isEditing = false;
	}

	function processWikiLinksForDisplay(html: string): string {
		const wikiLinkRegex = /\[\[([^\]]+)\]\]/g;
		return html.replace(wikiLinkRegex, (match, pageName) => {
			const trimmedName = pageName.trim();
			const linkedPage = $pagesStore.find((p) => p.name === trimmedName);
			
			if (linkedPage) {
				return `<a href="/page/${linkedPage.id}" class="wiki-link-exists" data-page-name="${trimmedName}">${trimmedName}</a>`;
			} else {
				return `<span class="wiki-link-missing cursor-pointer" data-page-name="${trimmedName}" onclick="window.handleMissingWikiLink('${trimmedName.replace(/'/g, "\\'")}')">${trimmedName}</span>`;
			}
		});
	}

	async function createNewPage(pageName: string) {
		creating = true;
		try {
			const newPage = await api.createPage({
				name: pageName,
				content: `<h1>${pageName}</h1><p>Start writing your notes here...</p>`
			});

			pagesStore.update((pages) => [...pages, newPage]);
			goto(`/page/${newPage.id}`);
			showCreateDialog = false;
			newPageName = '';
		} catch (error) {
			alert('Failed to create page: ' + (error as Error).message);
		} finally {
			creating = false;
		}
	}

	// Global function for handling missing wiki link clicks
	if (typeof window !== 'undefined') {
		(window as any).handleMissingWikiLink = (pageName: string) => {
			newPageName = pageName;
			showCreateDialog = true;
		};
	}

	// Watch for route changes
	$effect(() => {
		page = data.page;
		editedContent = data.page.content;
		isEditing = false;
		currentPageStore.set(data.page);
	});
</script>

<article class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-3xl font-bold">{page.name}</h1>
		<div class="flex gap-2">
			{#if !isEditing}
				<Button onclick={() => (isEditing = true)}>Edit</Button>
			{:else}
				<Button variant="outline" onclick={cancelEdit} disabled={isSaving}>Cancel</Button>
				<Button onclick={saveChanges} disabled={isSaving}>
					{isSaving ? 'Saving...' : 'Save'}
				</Button>
			{/if}
		</div>
	</div>

	{#if saveMessage}
		<div
			class="px-4 py-2 rounded-md {saveMessage.includes('Failed')
				? 'bg-red-100 text-red-800'
				: 'bg-green-100 text-green-800'}"
		>
			{saveMessage}
		</div>
	{/if}

	{#if isEditing}
		<Editor bind:content={editedContent} placeholder="Write your notes here..." />
	{:else}
		<div bind:this={contentElement} class="prose prose-sm lg:prose-lg max-w-none border border-gray-200 rounded-lg p-6">
			{@html displayContent}
		</div>
	{/if}

	<div class="text-xs text-gray-500 pt-4">
		Last updated: {new Date(page.updated_at).toLocaleString()}
	</div>
</article>

<Dialog bind:open={showCreateDialog}>
	<h2 class="text-xl font-bold mb-4">Create New Page</h2>
	<p class="mb-4 text-gray-600">
		The page "<strong>{newPageName}</strong>" doesn't exist yet. Would you like to create it?
	</p>
	<div class="flex gap-2 justify-end">
		<Button variant="outline" onclick={() => (showCreateDialog = false)}>Cancel</Button>
		<Button onclick={() => createNewPage(newPageName)} disabled={creating}>
			{creating ? 'Creating...' : 'Create Page'}
		</Button>
	</div>
</Dialog>

<style>
	:global(.wiki-link-exists) {
		color: #2563eb;
		text-decoration: underline;
		cursor: pointer;
	}

	:global(.wiki-link-missing) {
		color: #dc2626;
		text-decoration: underline;
		cursor: pointer;
	}

	:global(.wiki-link-exists:hover),
	:global(.wiki-link-missing:hover) {
		opacity: 0.8;
	}
</style>

