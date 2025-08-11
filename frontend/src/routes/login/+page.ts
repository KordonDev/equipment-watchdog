export async function load({ params, url }) {
	let username = url.searchParams.get('username');
	let message = url.searchParams.get('message');

	return {
		username,
		message
	};
}