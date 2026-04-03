import { getMember } from '$lib/services/member.service';
import { getOrdersForMember } from '$lib/services/order.service';
import { getNextGloveId } from '$lib/services/gloveId.service';
import type { PageLoad } from './$types';

export const prerender = false;

export const load: PageLoad = async ({ params }) => {
	const [member, orders, gloveIdResult] = await Promise.all([
		getMember(params.id),
		getOrdersForMember(params.id),
		getNextGloveId().catch(() => ({ nextId: null }))
	]);
	return { member, orders, nextGloveId: gloveIdResult.nextId };
};
