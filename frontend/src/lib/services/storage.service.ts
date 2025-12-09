import { Group, type Member } from '$lib/services/member.service';

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

const MEMBER_FOR_DETAILS = 'memberForDetails';
export const loadMemberForDetails = (members: Member[]): Member | undefined => {
	const stored = loadStoredValue(MEMBER_FOR_DETAILS);
	const member = members.find(m => m.id.toString() === stored);
	if (member) {
		return member;
	}
}
export function storeMemberForDetails(member: Member | undefined) {
	if (typeof window !== 'undefined') {
		if (member === undefined) {
			localStorage.removeItem(MEMBER_FOR_DETAILS);
			return;
		}
		localStorage.setItem(MEMBER_FOR_DETAILS, member.id.toString());
	}
}
