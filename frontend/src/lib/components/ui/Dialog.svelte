<script lang="ts">
	import { cn } from '$lib/utils/cn';

	let {
		open = $bindable(false),
		children
	}: {
		open?: boolean;
		children?: any;
	} = $props();

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			open = false;
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-50 bg-black/50 flex items-center justify-center"
		onclick={handleBackdropClick}
		role="presentation"
	>
		<div
			class={cn(
				'relative bg-white rounded-lg shadow-lg p-6 max-w-lg w-full mx-4',
				'animate-in fade-in-0 zoom-in-95'
			)}
		>
			{@render children?.()}
		</div>
	</div>
{/if}

