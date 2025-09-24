import { fetchApi } from '../apiService';


export function getNextGloveId(): Promise<GloveId> {
	return fetchApi('/glove-ids/next');
}

export function markGloveIdAsUsed(id: number): Promise<void> {
	return fetchApi(`/glove-ids/mark-used/${id}`, {
		method: 'POST'
	});
}

export interface GloveId {
	nextId: string;
}
