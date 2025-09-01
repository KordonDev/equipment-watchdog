<script lang="ts">
	import { onMount } from 'svelte';
	import { getAllUser, toggleApproveUser, toggleAdminUser, getMe, type User } from '$lib/services/user.service';

	let users: User[] = [];
	let loading = true;
	let currentUser: User | null = null;

	const loadUsers = async () => {
		loading = true;
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

	onMount(loadUsers);
</script>

<div class="p-6">
	<h1 class="text-2xl font-bold mb-6">Benutzerverwaltung</h1>
	{#if loading}
		<div>Lade Benutzer...</div>
	{:else if users.length === 0}
		<div>Keine Benutzer gefunden.</div>
	{:else}
		<table class="min-w-full bg-white border border-gray-200 rounded">
			<thead>
				<tr>
					<th class="px-4 py-2 border-b">ID</th>
					<th class="px-4 py-2 border-b">Name</th>
					<th class="px-4 py-2 border-b">Genehmigt</th>
					<th class="px-4 py-2 border-b">Admin</th>
					<th class="px-4 py-2 border-b">Aktionen</th>
				</tr>
			</thead>
			<tbody>
				{#each users as user}
					<tr>
						<td class="px-4 py-2 border-b">{user.id}</td>
						<td class="px-4 py-2 border-b">{user.name}</td>
						<td class="px-4 py-2 border-b">
							{user.isApproved ? '✅' : '❌'}
						</td>
						<td class="px-4 py-2 border-b">
							{user.isAdmin ? '✅' : '❌'}
						</td>
						<td class="px-4 py-2 border-b space-x-2">
							{#if currentUser && currentUser.isAdmin}
								<button
									class="px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
									on:click={() => handleToggleApprove(user)}
								>
									{user.isApproved ? 'Entfernen' : 'Genehmigen'}
								</button>
								<button
									class="px-2 py-1 bg-purple-600 text-white rounded hover:bg-purple-700"
									on:click={() => handleToggleAdmin(user)}
								>
									{user.isAdmin ? 'Admin entfernen' : 'Zum Admin machen'}
								</button>
							{/if}
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	{/if}
</div>
