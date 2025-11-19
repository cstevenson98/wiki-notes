<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Link from '@tiptap/extension-link';
	import Placeholder from '@tiptap/extension-placeholder';
	import Mathematics from '@tiptap/extension-mathematics';
	import 'katex/dist/katex.min.css';
	import { pagesStore } from '$lib/stores/pages';
	import { goto } from '$app/navigation';
	import { api } from '$lib/utils/api';
	import Dialog from './ui/Dialog.svelte';
	import Input from './ui/Input.svelte';
	import Button from './ui/Button.svelte';

	let {
		content = $bindable(''),
		placeholder = 'Start writing...'
	}: {
		content?: string;
		placeholder?: string;
	} = $props();

	let editorElement: HTMLDivElement;
	let editor: Editor | null = null;
	let showCreateDialog = $state(false);
	let newPageName = $state('');
	let creating = $state(false);
	let wikiLinkTimeout: ReturnType<typeof setTimeout>;

	onMount(() => {
		editor = new Editor({
			element: editorElement,
			extensions: [
				StarterKit,
				Link.configure({
					openOnClick: false,
					HTMLAttributes: {
						class: 'text-blue-600 underline'
					}
				}),
				Placeholder.configure({
					placeholder
				}),
			Mathematics
			],
			content,
			editorProps: {
				attributes: {
					class:
						'prose prose-sm sm:prose lg:prose-lg xl:prose-xl focus:outline-none max-w-none min-h-[400px] p-4'
				},
				handleClickOn: (view, pos, node, nodePos, event) => {
					// Check if clicking on a wiki link
					const target = event.target as HTMLElement;
					if (target.classList.contains('wiki-link-exists')) {
						const pageName = target.getAttribute('data-page-name');
						if (pageName) {
							handleWikiLinkClick(pageName, true);
						}
						return true;
					} else if (target.classList.contains('wiki-link-missing')) {
						const pageName = target.getAttribute('data-page-name');
						if (pageName) {
							handleWikiLinkClick(pageName, false);
						}
						return true;
					}
					return false;
				}
			},
			onUpdate: ({ editor }) => {
				content = editor.getHTML();
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});

	function processWikiLinks() {
		if (!editor) return;

		const html = editor.getHTML();
		const wikiLinkRegex = /\[\[([^\]]+)\]\]/g;
		let processedHtml = html;

		const matches = Array.from(html.matchAll(wikiLinkRegex));
		for (const match of matches) {
			const fullMatch = match[0];
			const pageName = match[1].trim();
			const page = $pagesStore.find((p) => p.name === pageName);

			const replacement = `<span data-type="wikilink" data-page-name="${pageName}" data-exists="${!!page}" class="${page ? 'wiki-link-exists' : 'wiki-link-missing'}">${fullMatch}</span>`;
			processedHtml = processedHtml.replace(fullMatch, replacement);
		}

		if (processedHtml !== html) {
			const currentPos = editor.state.selection.anchor;
			editor.commands.setContent(processedHtml);
			editor.commands.setTextSelection(currentPos);
		}
	}

	async function handleWikiLinkClick(pageName: string, exists: boolean) {
		if (exists) {
			const page = $pagesStore.find((p) => p.name === pageName);
			if (page) {
				goto(`/page/${page.id}`);
			}
		} else {
			newPageName = pageName;
			showCreateDialog = true;
		}
	}

	async function createPage() {
		if (!newPageName.trim()) return;

		creating = true;
		try {
			const newPage = await api.createPage({
				name: newPageName,
				content: `# ${newPageName}\n\nStart writing your notes here...`
			});

			// Update pages store
			pagesStore.update((pages) => [...pages, newPage]);

			// Navigate to new page
			goto(`/page/${newPage.id}`);
			showCreateDialog = false;
			newPageName = '';
		} catch (error) {
			alert('Failed to create page: ' + (error as Error).message);
		} finally {
			creating = false;
		}
	}

	// Watch for content changes from parent
	$effect(() => {
		if (editor && content !== editor.getHTML()) {
			editor.commands.setContent(content);
			processWikiLinks();
		}
	});

	// Re-process wiki links when pages change
	$effect(() => {
		if ($pagesStore) {
			processWikiLinks();
		}
	});
</script>

<div class="border border-gray-300 rounded-lg overflow-hidden">
	<div bind:this={editorElement}></div>
</div>

<Dialog bind:open={showCreateDialog}>
	<h2 class="text-xl font-bold mb-4">Create New Page</h2>
	<p class="mb-4 text-gray-600">
		The page "<strong>{newPageName}</strong>" doesn't exist yet. Would you like to create it?
	</p>
	<div class="flex gap-2 justify-end">
		<Button variant="outline" onclick={() => (showCreateDialog = false)}>Cancel</Button>
		<Button onclick={createPage} disabled={creating}>
			{creating ? 'Creating...' : 'Create Page'}
		</Button>
	</div>
</Dialog>

<style>
	:global(.wiki-link-exists) {
		color: #2563eb;
		cursor: pointer;
		text-decoration: underline;
	}

	:global(.wiki-link-missing) {
		color: #dc2626;
		cursor: pointer;
		text-decoration: underline;
	}

	:global(.wiki-link-exists:hover),
	:global(.wiki-link-missing:hover) {
		opacity: 0.8;
	}

	:global(.math-node) {
		display: inline-block;
		padding: 0 4px;
	}
</style>

