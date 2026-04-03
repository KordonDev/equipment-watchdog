import { Group } from '$lib/services/member.service';

const loadStoredValue = (storageKey: string): string | undefined => {
	if (typeof window !== 'undefined') {
		return localStorage.getItem(storageKey) || undefined;
	}
};

const SELECTED_GROUP_KEY = 'selectedGroup';
export const loadSelectedGroup = (): Group => {
	const stored = loadStoredValue(SELECTED_GROUP_KEY);
	if (stored && Object.values(Group).includes(stored as Group)) {
		return stored as Group;
	}
	return Group.FRIDAY;
}
export function storeSelectedGroup(selectedGroup: Group) {
	if (typeof window !== 'undefined') {
		localStorage.setItem(SELECTED_GROUP_KEY, selectedGroup);
	}
}
