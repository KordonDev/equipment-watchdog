import { fetchApi } from '../apiService';


export function getNextGloveId(): Promise<GloveId> {
	return fetchApi('/glove-ids/next');
}

export interface GloveId {
	nextId: string;
}
