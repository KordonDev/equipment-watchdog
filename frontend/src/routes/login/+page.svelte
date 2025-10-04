<script lang="ts">
	import type { PageProps } from './$types';
	import { login, setStoredUsername, getStoredUsername } from "$lib/services/authentication";
	import { goto } from '$app/navigation';

	let { data }: PageProps = $props();

	let infoMessage: string | null = $state(data.message);
	let username = $state(data.username || '');
	let storeUsername = $state(false);
	let loading = $state(false);

	// If username is not set from URL, try to get it from localStorage
	if (!data.username) {
		const stored = getStoredUsername();
		if (stored) {
			username = stored;
		}
	}

	const handleLogin = async (event: SubmitEvent) => {
		event.preventDefault();
		loading = true;
		infoMessage = null;
		try {
			await login(username);
			if (storeUsername) {
				setStoredUsername(username);
			} else {
				setStoredUsername('');
			}
			goto('/members');
		} catch (e) {
			infoMessage = "Login fehlgeschlagen. Bitte überprüfe deinen Benutzernamen.";
		} finally {
			loading = false;
		}
	};
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
	<div class="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
		<h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Login</h2>
		{#if infoMessage}
			<div class="mb-4 text-sm text-green-700 bg-green-100 rounded px-3 py-2">{infoMessage}</div>
		{/if}
		<form onsubmit={handleLogin} class="space-y-4">
			<div>
				<label for="username" class="block mb-2 text-sm font-medium text-gray-700">Benutzer:</label>
				<input
					required
					class="block w-full rounded border border-gray-300 px-3 py-2 focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
					id="username"
					bind:value={username}
					type="text"
					autocomplete="username webauthn"
				/>
			</div>
			<div class="flex items-center">
				<input
					type="checkbox"
					bind:checked={storeUsername}
					id="storeUsername"
					class="mr-2 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
				/>
				<label for="storeUsername" class="text-sm text-gray-700 select-none">Benutzername speichern</label>
			</div>
			<button
				type="submit"
				class="w-full rounded bg-purple-600 px-4 py-2 text-white font-semibold hover:bg-purple-700 transition-colors duration-150 disabled:opacity-50"
				disabled={loading}
			>
				Einloggen
			</button>
		</form>
	</div>
</div>
