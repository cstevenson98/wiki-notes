import { Mark, mergeAttributes } from '@tiptap/core';

export interface WikiLinkOptions {
	HTMLAttributes: Record<string, any>;
	onLinkClick?: (pageName: string) => void;
}

export const WikiLink = Mark.create<WikiLinkOptions>({
	name: 'wikilink',

	priority: 1000,

	addOptions() {
		return {
			HTMLAttributes: {},
			onLinkClick: undefined
		};
	},

	addAttributes() {
		return {
			pageName: {
				default: null,
				parseHTML: (element) => element.getAttribute('data-page-name'),
				renderHTML: (attributes) => {
					return {
						'data-page-name': attributes.pageName
					};
				}
			},
			exists: {
				default: true,
				parseHTML: (element) => element.getAttribute('data-exists') === 'true',
				renderHTML: (attributes) => {
					return {
						'data-exists': attributes.exists
					};
				}
			}
		};
	},

	parseHTML() {
		return [
			{
				tag: 'span[data-type="wikilink"]'
			}
		];
	},

	renderHTML({ HTMLAttributes }) {
		return [
			'span',
			mergeAttributes(
				{
					'data-type': 'wikilink',
					class: HTMLAttributes.exists === 'true' ? 'wiki-link-exists' : 'wiki-link-missing'
				},
				this.options.HTMLAttributes,
				HTMLAttributes
			),
			0
		];
	},

	addCommands() {
		return {
			setWikiLink:
				(attributes) =>
				({ commands }) => {
					return commands.setMark(this.name, attributes);
				},
			unsetWikiLink:
				() =>
				({ commands }) => {
					return commands.unsetMark(this.name);
				}
		};
	}
});

