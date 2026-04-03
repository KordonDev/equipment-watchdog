import { getMember } from '$lib/services/member.service';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ params }) => {
	const member = await getMember(params.id);
	return { member };
};
