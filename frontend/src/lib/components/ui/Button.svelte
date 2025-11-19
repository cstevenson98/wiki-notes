<script lang="ts">
	import { cn } from '$lib/utils/cn';

	let {
		variant = 'default',
		size = 'default',
		class: className,
		children,
		...restProps
	}: {
		variant?: 'default' | 'outline' | 'ghost' | 'destructive';
		size?: 'default' | 'sm' | 'lg';
		class?: string;
		children?: any;
		[key: string]: any;
	} = $props();

	const variantClasses = {
		default: 'bg-black text-white hover:bg-gray-800',
		outline: 'border border-black bg-white hover:bg-gray-100',
		ghost: 'hover:bg-gray-100',
		destructive: 'bg-red-600 text-white hover:bg-red-700'
	};

	const sizeClasses = {
		default: 'h-10 px-4 py-2',
		sm: 'h-9 px-3 text-sm',
		lg: 'h-11 px-8'
	};
</script>

<button
	class={cn(
		'inline-flex items-center justify-center rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
		variantClasses[variant],
		sizeClasses[size],
		className
	)}
	{...restProps}
>
	{@render children?.()}
</button>

