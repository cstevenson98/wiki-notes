import { writable } from 'svelte/store';
import type { Page } from '$lib/utils/api';

export const pagesStore = writable<Page[]>([]);
export const currentPageStore = writable<Page | null>(null);

