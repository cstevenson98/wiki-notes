import { redirect } from '@sveltejs/kit';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	// Use relative URL - SvelteKit's fetch will handle it correctly
	const response = await fetch('/api/page/by-name/Home');
	if (response.ok) {
		const page = await response.json();
		throw redirect(302, `/page/${page.id}`);
	}

	// If Home page doesn't exist, stay on this page
	return {};
};

