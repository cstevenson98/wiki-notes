<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import Link from '@tiptap/extension-link';
	import Placeholder from '@tiptap/extension-placeholder';
	import Mathematics from '@tiptap/extension-mathematics';
	import EditorToolbar from './EditorToolbar.svelte';
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
	<EditorToolbar {editor} updateKey={toolbarKey} />

	<!-- Editor -->
	<div bind:this={editorElement}></div>
</div>


<style>
	:global(.math-node) {
		display: inline-block;
		padding: 0 4px;
	}
</style>

