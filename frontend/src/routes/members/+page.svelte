<script lang="ts">
	import { onMount } from 'svelte';
	import { getMembers, Group, type Member, createMember } from '$lib/services/member.service';
	import { EquipmentType } from '$lib/services/equipment.service';

	let members: Member[] = $state([]);
	let selectedGroup: Group = $state(Group.FRIDAY);
	let loading = $state(true);
	let showDialog = $state(false);
	let newMemberName = $state('');
	let addingMember = $state(false);

	let filteredMembers = $derived(members.filter(member => member.group === selectedGroup));

	onMount(async () => {
		try {
			members = await getMembers();
		} catch (error) {
			console.error('Failed to load members:', error);
		} finally {
			loading = false;
		}
	});

	const groupLabels = {
		[Group.FRIDAY]: 'Freitag',
		[Group.MONDAY]: 'Montag',
		[Group.MINI]: 'Mini'
	};

	const equipmentLabels = {
		[EquipmentType.Helmet]: 'Helm',
		[EquipmentType.Jacket]: 'Jacke',
		[EquipmentType.Gloves]: 'Handschuhe',
		[EquipmentType.Trousers]: 'Hose',
		[EquipmentType.Boots]: 'Stiefel',
		[EquipmentType.TShirt]: 'T-Shirt'
	};

	const handleAddMember = () => {
		showDialog = true;
		newMemberName = '';
	};

	const handleSubmitMember = async () => {
		if (!newMemberName.trim()) return;

		addingMember = true;
		try {
			const newMember: Member = {
				id: '', // Will be set by backend
				name: newMemberName.trim(),
				group: selectedGroup,
				equipments: {}
			};

			const createdMember = await createMember(newMember);
			members = [...members, createdMember];
			showDialog = false;
		} catch (error) {
			console.error('Failed to create member:', error);
		} finally {
			addingMember = false;
		}
	};

	const handleCancelDialog = () => {
		showDialog = false;
		newMemberName = '';
	};
</script>

<div class="p-6">
	<h1 class="text-2xl font-bold mb-6">Personen</h1>

	<!-- Group Switcher with Add Button -->
	<div class="mb-6 flex items-center justify-between">
		<div class="flex space-x-2">
			{#each Object.values(Group) as group}
				<button
					class="px-4 py-2 rounded {selectedGroup === group ? 'bg-purple-600 text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'}"
					onclick={() => selectedGroup = group}
				>
					{groupLabels[group]}
				</button>
			{/each}
		</div>
		<button
			class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
			onclick={handleAddMember}
		>
			+ Mitglied hinzufügen
		</button>
	</div>

	<!-- Members List -->
	{#if loading}
		<div class="text-center">Laden...</div>
	{:else if filteredMembers.length === 0}
		<div class="text-center text-gray-500">Keine Mitglieder in dieser Gruppe gefunden.</div>
	{:else}
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
			{#each filteredMembers as member}
				<div class="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow">
					<h3 class="text-lg font-semibold text-gray-800 mb-3">{member.name}</h3>

					<!-- Equipment Indicators -->
					<div class="flex flex-wrap gap-1">
						{#each Object.values(EquipmentType) as equipmentType}
							{@const hasEquipment = member.equipments[equipmentType]}
							<span
								class="px-2 py-1 text-xs rounded flex items-center gap-1 {hasEquipment ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}"
								title="{equipmentLabels[equipmentType]}: {hasEquipment ? 'Vorhanden' : 'Fehlt'}"
							>
								<span class="font-bold">{hasEquipment ? '✓' : '✗'}</span>
								{equipmentLabels[equipmentType]}
							</span>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	{/if}
</div>

<!-- Add Member Dialog -->
{#if showDialog}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-md">
			<h2 class="text-xl font-bold mb-4">Neues Mitglied hinzufügen</h2>
			<p class="text-sm text-gray-600 mb-4">Gruppe: {groupLabels[selectedGroup]}</p>

			<div class="mb-4">
				<label for="memberName" class="block mb-2 text-sm font-medium text-gray-700">Name:</label>
				<input
					id="memberName"
					type="text"
					bind:value={newMemberName}
					class="block w-full rounded border border-gray-300 px-3 py-2 focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
					placeholder="Mitgliedname eingeben..."
					onkeydown={(e) => e.key === 'Enter' && handleSubmitMember()}
				/>
			</div>

			<div class="flex justify-end space-x-2">
				<button
					type="button"
					onclick={handleCancelDialog}
					class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
					disabled={addingMember}
				>
					Abbrechen
				</button>
				<button
					type="button"
					onclick={handleSubmitMember}
					class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors disabled:opacity-50"
					disabled={addingMember || !newMemberName.trim()}
				>
					{addingMember ? 'Hinzufügen...' : 'Hinzufügen'}
				</button>
			</div>
		</div>
	</div>
{/if}
