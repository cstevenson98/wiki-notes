<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Link from '@tiptap/extension-link';
	import Placeholder from '@tiptap/extension-placeholder';
	import Mathematics from '@tiptap/extension-mathematics';
	import 'katex/dist/katex.min.css';

	let {
		content = $bindable(''),
		placeholder = 'Start writing...'
	}: {
		content?: string;
		placeholder?: string;
	} = $props();

	let editorElement: HTMLDivElement;
	let editor: Editor | null = null;
	
	// Toolbar state - only updated on selection changes, not every keystroke
	let toolbarKey = $state(0);
	let isSettingContent = false;

	onMount(() => {
		editor = new Editor({
			element: editorElement,
			extensions: [
				StarterKit.configure({
					// Exclude Link from StarterKit so we can configure it ourselves
					link: false
				}),
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
				handleClickOn: () => {
					return false;
				}
			},
			onUpdate: ({ editor }) => {
				if (!isSettingContent) {
					content = editor.getHTML();
				}
			},
			onSelectionUpdate: () => {
				// Increment key to trigger toolbar re-render only on selection changes
				toolbarKey++;
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});



	// Watch for content changes from parent - only update editor if content prop changes
	let lastContentFromProp = content;
	$effect(() => {
		if (editor && content !== lastContentFromProp && content !== editor.getHTML()) {
			lastContentFromProp = content;
			isSettingContent = true;
			editor.commands.setContent(content);
			isSettingContent = false;
		}
	});
</script>

<div class="border border-gray-300 rounded-lg overflow-hidden">
	<!-- Toolbar - only re-renders on selection changes via toolbarKey -->
	{#key toolbarKey}
	<div class="border-b border-gray-300 bg-gray-50 p-2 flex flex-wrap gap-1">
		<!-- Text formatting -->
		<button
			onclick={() => editor?.chain().focus().toggleBold().run()}
			class="toolbar-btn {editor?.isActive('bold') ? 'active' : ''}"
			title="Bold (Ctrl+B)"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 4h8a4 4 0 014 4 4 4 0 01-4 4H6z M6 12h9a4 4 0 014 4 4 4 0 01-4 4H6z"></path>
			</svg>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleItalic().run()}
			class="toolbar-btn {editor?.isActive('italic') ? 'active' : ''}"
			title="Italic (Ctrl+I)"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 4h4 M14 4l-4 16 M10 20h4"></path>
			</svg>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleStrike().run()}
			class="toolbar-btn {editor?.isActive('strike') ? 'active' : ''}"
			title="Strikethrough"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12h18 M8 5a4 4 0 018 0 M8 19a4 4 0 008 0"></path>
			</svg>
		</button>

		<button
			onclick={() => editor?.chain().focus().toggleCode().run()}
			class="toolbar-btn {editor?.isActive('code') ? 'active' : ''}"
			title="Inline Code"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4"></path>
			</svg>
		</button>

		<div class="w-px h-6 bg-gray-300 mx-1"></div>

		<!-- Headings -->
		<button
			onclick={() => editor?.chain().focus().toggleHeading({ level: 1 }).run()}
			class="toolbar-btn {editor?.isActive('heading', { level: 1 }) ? 'active' : ''}"
			title="Heading 1"
			type="button"
		>
			<span class="font-bold">H1</span>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleHeading({ level: 2 }).run()}
			class="toolbar-btn {editor?.isActive('heading', { level: 2 }) ? 'active' : ''}"
			title="Heading 2"
			type="button"
		>
			<span class="font-bold">H2</span>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleHeading({ level: 3 }).run()}
			class="toolbar-btn {editor?.isActive('heading', { level: 3 }) ? 'active' : ''}"
			title="Heading 3"
			type="button"
		>
			<span class="font-bold">H3</span>
		</button>

		<div class="w-px h-6 bg-gray-300 mx-1"></div>

		<!-- Lists -->
		<button
			onclick={() => editor?.chain().focus().toggleBulletList().run()}
			class="toolbar-btn {editor?.isActive('bulletList') ? 'active' : ''}"
			title="Bullet List"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"></path>
			</svg>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleOrderedList().run()}
			class="toolbar-btn {editor?.isActive('orderedList') ? 'active' : ''}"
			title="Numbered List"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h1v5H3 M3 14h1v5H3 M8 6h13M8 12h13M8 18h13"></path>
			</svg>
		</button>

		<div class="w-px h-6 bg-gray-300 mx-1"></div>

		<!-- Blocks -->
		<button
			onclick={() => editor?.chain().focus().toggleCodeBlock().run()}
			class="toolbar-btn {editor?.isActive('codeBlock') ? 'active' : ''}"
			title="Code Block"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
			</svg>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().toggleBlockquote().run()}
			class="toolbar-btn {editor?.isActive('blockquote') ? 'active' : ''}"
			title="Blockquote"
			type="button"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
			</svg>
		</button>

		<div class="w-px h-6 bg-gray-300 mx-1"></div>

		<!-- Undo/Redo -->
		<button
			onclick={() => editor?.chain().focus().undo().run()}
			class="toolbar-btn"
			title="Undo (Ctrl+Z)"
			type="button"
			disabled={!editor?.can().undo()}
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6"></path>
			</svg>
		</button>
		
		<button
			onclick={() => editor?.chain().focus().redo().run()}
			class="toolbar-btn"
			title="Redo (Ctrl+Shift+Z)"
			type="button"
			disabled={!editor?.can().redo()}
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10h-10a8 8 0 00-8 8v2M21 10l-6 6m6-6l-6-6"></path>
			</svg>
		</button>
	</div>
	{/key}

	<!-- Editor -->
	<div bind:this={editorElement}></div>
</div>


<style>
	:global(.math-node) {
		display: inline-block;
		padding: 0 4px;
	}

	.toolbar-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 0.375rem 0.5rem;
		border-radius: 0.375rem;
		background-color: white;
		border: 1px solid #e5e7eb;
		color: #374151;
		font-size: 0.875rem;
		cursor: pointer;
		transition: all 0.15s;
		min-width: 2rem;
	}

	.toolbar-btn:hover:not(:disabled) {
		background-color: #f3f4f6;
		border-color: #d1d5db;
	}

	.toolbar-btn:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}

	.toolbar-btn.active {
		background-color: #3b82f6;
		color: white;
		border-color: #2563eb;
	}

	.toolbar-btn.active:hover {
		background-color: #2563eb;
	}
</style>

