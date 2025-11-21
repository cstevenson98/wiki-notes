const API_URL = import.meta.env.PUBLIC_API_URL || 'http://localhost:8080';

export interface Page {
	id: number;
	name: string;
	content: string;
	created_at: string;
	updated_at: string;
}

export interface CreatePageRequest {
	name: string;
	content: string;
}

export interface UpdatePageRequest {
	name?: string;
	content?: string;
}

export const api = {
	async getAllPages(): Promise<Page[]> {
		const response = await fetch(`${API_URL}/api/pages`);
		if (!response.ok) {
			throw new Error('Failed to fetch pages');
		}
		return response.json();
	},

	async getPageById(id: number): Promise<Page> {
		const response = await fetch(`${API_URL}/api/page/${id}`);
		if (!response.ok) {
			throw new Error('Failed to fetch page');
		}
		return response.json();
	},

	async getPageByName(name: string): Promise<Page | null> {
		const response = await fetch(`${API_URL}/api/page/by-name/${encodeURIComponent(name)}`);
		if (response.status === 404) {
			return null;
		}
		if (!response.ok) {
			throw new Error('Failed to fetch page');
		}
		return response.json();
	},

	async createPage(data: CreatePageRequest): Promise<Page> {
		const response = await fetch(`${API_URL}/api/page`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});
		if (!response.ok) {
			const error = await response.json();
			throw new Error(error.error || 'Failed to create page');
		}
		return response.json();
	},

	async updatePage(id: number, data: UpdatePageRequest): Promise<Page> {
		const response = await fetch(`${API_URL}/api/page/${id}`, {
			method: 'PATCH',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});
		if (!response.ok) {
			throw new Error('Failed to update page');
		}
		return response.json();
	},

	async deletePage(id: number): Promise<void> {
		const response = await fetch(`${API_URL}/api/page/${id}`, {
			method: 'DELETE'
		});
		if (!response.ok) {
			throw new Error('Failed to delete page');
		}
	},

	async getBacklinks(id: number): Promise<Page[]> {
		const response = await fetch(`${API_URL}/api/page/${id}/backlinks`);
		if (!response.ok) {
			throw new Error('Failed to fetch backlinks');
		}
		const data = await response.json();
		// Ensure we always return an array, even if the response is null
		return Array.isArray(data) ? data : [];
	}
};

