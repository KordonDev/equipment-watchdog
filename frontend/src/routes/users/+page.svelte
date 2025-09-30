<script lang="ts">
	import { onMount } from 'svelte';
	import { getAllUser, toggleApproveUser, toggleAdminUser, getMe, type User } from '$lib/services/user.service';
	import { changePassword, addLogin } from '$lib/services/authentication';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';

	let users: User[] = [];
	let loading = true;
	let currentUser: User | null = null;
	let newPassword = '';
	let passwordChangeLoading = false;
	let passwordChangeSuccess = false;
	let passwordChangeError = '';
	let passkeyLoading = false;
	let passkeySuccess = false;
	let passkeyError = '';

	const loadUsers = async () => {
		currentUser = await getMe();
		if (currentUser.isAdmin) {
			users = await getAllUser();
		} else {
			users = [currentUser];
		}
		loading = false;
	};

	const handleToggleApprove = async (user: User) => {
		await toggleApproveUser(user.name);
		await loadUsers();
	};

	const handleToggleAdmin = async (user: User) => {
		await toggleAdminUser(user.name);
		await loadUsers();
	};

	const handleChangePassword = async () => {
		passwordChangeLoading = true;
		passwordChangeSuccess = false;
		passwordChangeError = '';
		try {
			await changePassword(newPassword);
			passwordChangeSuccess = true;
			newPassword = '';
		} catch (e) {
			passwordChangeError = 'Fehler beim Ändern des Passworts.';
		}
		passwordChangeLoading = false;
	};

	const handleAddPasskey = async () => {
		passkeyLoading = true;
		passkeySuccess = false;
		passkeyError = '';
		try {
			await addLogin();
			passkeySuccess = true;
		} catch (e) {
			console.error(e)
			passkeyError = 'Fehler beim Hinzufügen des Passkeys.';
		}
		passkeyLoading = false;
	};

	onMount(loadUsers);
</script>

<div class="p-6">
	<BurgerMenu />
	<h1 class="text-2xl font-bold mb-6">Benutzer</h1>
	{#if loading}
		<div>Lade Benutzer...</div>
	{:else if users.length === 0}
		<div>Keine Benutzer gefunden.</div>
	{:else}
		<table class="min-w-full border-collapse border border-gray-200">
			<thead>
				<tr>
					<th class="px-4 py-2 border-b text-left">Name</th>
					<th class="px-4 py-2 border-b text-left">Genehmigt</th>
					<th class="px-4 py-2 border-b text-left">Admin</th>
					<th class="px-4 py-2 border-b text-left">Aktionen</th>
				</tr>
			</thead>
			<tbody>
				{#each users as user (user.id)}
					<tr>
						<td class="px-4 py-2 border-b">{user.name}</td>
						<td class="px-4 py-2 border-b">
							{user.isApproved ? '✅' : '❌'}
						</td>
						<td class="px-4 py-2 border-b">
							{user.isAdmin ? '✅' : '❌'}
						</td>
						<td class="px-4 py-2 border-b space-x-2">
							{#if currentUser && currentUser.isAdmin && currentUser.id !== user.id}
								<button
									class="px-2 py-1 mb-2 bg-blue-600 text-white rounded hover:bg-blue-700"
									on:click={() => handleToggleApprove(user)}
								>
									{user.isApproved ? 'Entfernen' : 'Genehmigen'}
								</button>
								{#if currentUser && user.isApproved}
									<button
										class="px-2 py-1 bg-purple-600 text-white rounded hover:bg-purple-700"
										on:click={() => handleToggleAdmin(user)}
									>
										{user.isAdmin ? 'Admin entfernen' : 'Admin machen'}
									</button>
								{/if}
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}

	<!-- Passwort ändern Bereich -->
	{#if currentUser}
		<div class="mt-8 max-w-md">
			<h2 class="text-lg font-semibold mb-2">Passwort ändern</h2>
			<div class="flex gap-2 items-center">
				<input
					type="password"
					placeholder="Neues Passwort"
					class="px-2 py-1 border rounded flex-1"
					bind:value={newPassword}
				/>
				<button
					class="px-3 py-1 bg-green-600 text-white rounded hover:bg-green-700"
					disabled={passwordChangeLoading || !newPassword}
					on:click={handleChangePassword}
				>
					Ändern
				</button>
			</div>
			{#if passwordChangeSuccess}
				<div class="text-green-600 mt-2">Passwort erfolgreich geändert.</div>
			{/if}
			{#if passwordChangeError}
				<div class="text-red-600 mt-2">{passwordChangeError}</div>
			{/if}

			<h2 class="text-lg font-semibold mt-8 mb-2">Neuen Passkey hinzufügen</h2>
			<div class="flex gap-2 items-center">
				<button
					class="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
					disabled={passkeyLoading}
					on:click={handleAddPasskey}
				>
					Passkey hinzufügen
				</button>
			</div>
			{#if passkeySuccess}
				<div class="text-green-600 mt-2">Passkey erfolgreich hinzugefügt.</div>
			{/if}
			{#if passkeyError}
				<div class="text-red-600 mt-2">{passkeyError}</div>
			{/if}
		</div>
	{/if}
</div>
