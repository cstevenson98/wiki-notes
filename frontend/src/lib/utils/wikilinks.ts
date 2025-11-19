import type { Page } from './api';

export interface ParsedWikiLink {
	text: string;
	pageName: string;
	exists: boolean;
	pageId?: string;
}

export function parseWikiLinks(content: string, pages: Page[]): ParsedWikiLink[] {
	const wikiLinkRegex = /\[\[([^\]]+)\]\]/g;
	const links: ParsedWikiLink[] = [];
	let match;

	while ((match = wikiLinkRegex.exec(content)) !== null) {
		const pageName = match[1].trim();
		const page = pages.find((p) => p.name === pageName);

		links.push({
			text: match[0],
			pageName,
			exists: !!page,
			pageId: page?.id
		});
	}

	return links;
}

export function extractWikiLinkNames(content: string): string[] {
	const wikiLinkRegex = /\[\[([^\]]+)\]\]/g;
	const names: string[] = [];
	let match;

	while ((match = wikiLinkRegex.exec(content)) !== null) {
		names.push(match[1].trim());
	}

	return names;
}

