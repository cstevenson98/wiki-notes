import { error } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ params, fetch }) => {
	try {
		// Use relative URL - SvelteKit's fetch will handle it correctly
		const response = await fetch(`/api/page/${params.id}`);

		if (!response.ok) {
			if (response.status === 404) {
				throw error(404, 'Page not found');
			}
			throw error(500, 'Failed to load page');
		}

		const page = await response.json();
		return { page };
	} catch (err) {
		if (err instanceof Response) {
			throw err;
		}
		console.error('Failed to load page:', err);
		throw error(500, 'Failed to load page');
	}
};

