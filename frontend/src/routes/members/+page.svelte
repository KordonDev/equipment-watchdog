<script lang="ts">
	import { onMount } from 'svelte';
	import { equipmentForGroup, getMembers, Group, groupLabels, type Member } from '$lib/services/member.service';
	import { equipmentLabels, EquipmentType } from '$lib/services/equipment.service';
	import { getOrders, type Order } from '$lib/services/order.service';
	import {
		loadMemberForDetails,
		loadSelectedGroup,
		storeMemberForDetails,
		storeSelectedGroup
	} from '$lib/services/storage.service';
	import AddMemberDialog from '$lib/components/AddMemberDialog.svelte';
	import MemberDetailDialog from '$lib/components/MemberDetailDialog.svelte';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';

	let members: Member[] = $state([]);
	let selectedGroup: Group = $state(loadSelectedGroup());
	let loading = $state(true);
	let showAddMemberDialog = $state(false);
	let memberForDetails: Member | undefined = $state(undefined);
	let orders: Order[] = $state([]);

	let filteredMembers = $derived(members.filter(member => member.group === selectedGroup));

	onMount(async () => {
		try {
			members = await getMembers();
			orders = await getOrders();
			memberForDetails = loadMemberForDetails(members)
		} catch (error) {
			console.error('Failed to load members or orders:', error);
		} finally {
			loading = false;
		}
	});

	$effect(() => {
		storeSelectedGroup(selectedGroup)
	});

	const handleAddMember = () => {
		showAddMemberDialog = true;
	};

	const handleMemberAdded = (newMember: Member) => {
		members = [...members, newMember];
		showAddMemberDialog = false;
	};

	const handleCloseDialog = () => {
		showAddMemberDialog = false;
	};

	const handleMemberClick = (member: Member) => {
		memberForDetails = member;
		storeMemberForDetails(member)
	};

	const handleMemberUpdated = async (updatedMember: Member) => {
		const index = members.findIndex(m => m.id === updatedMember.id);
		if (index !== -1) {
			members[index] = updatedMember;
			members = [...members]; // Trigger reactivity and update filteredMembers
		}
		// Nach Equipment-Änderung Bestellungen neu laden
		orders = await getOrders();
	};

	const handleCloseDetailDialog = (deletedMemberId?: string) => {
		if (deletedMemberId) {
			members = members.filter(m => m.id.toString() !== deletedMemberId);
		}
		memberForDetails = undefined;
		storeMemberForDetails(undefined)
		members = [...members]; // Trigger reactivity
	};

	const hasOrderForMemberAndType = (memberId: number, type: EquipmentType) =>
		orders.some(order => order.memberId === memberId && order.type === type && !order.fulfilledAt);
</script>

<div class="p-6">
	<BurgerMenu />
	<h1 class="text-2xl font-bold mb-6">Mitglieder</h1>

	<!-- Group Switcher with Add Button -->
	<div class="mb-6 flex items-center justify-between flex-wrap gap-2">
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
			class="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 transition-colors ml-auto"
			onclick={handleAddMember}
		>
			+ Mitglied
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
				<button
					class="bg-white rounded-lg shadow-md p-4 hover:shadow-lg transition-shadow cursor-pointer"
					onclick={() => handleMemberClick(member)}
				>
					<h3 class="text-lg font-semibold text-gray-800 mb-3">{member.name}</h3>

					<!-- Equipment Indicators -->
					<div class="flex flex-wrap gap-1">
						{#each (equipmentForGroup[member.group] || []) as equipmentType}
							{@const hasEquipment = member.equipments[equipmentType]}
							<span
								class="px-2 py-1 text-xs rounded flex items-center gap-1 {hasEquipment ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}"
								title="{equipmentLabels[equipmentType]}: {hasEquipment ? 'Vorhanden' : 'Fehlt'}"
							>
								<span class="font-bold">{hasEquipment ? '✓' : '✗'}</span>
								{equipmentLabels[equipmentType]}
								{#if hasOrderForMemberAndType(member.id, equipmentType)}
									<span class="ml-1 px-1 py-0.5 rounded bg-yellow-100 text-yellow-800" title="Bestellt">🛒</span>
								{/if}
							</span>
						{/each}
					</div>
				</button>
			{/each}
		</div>
	{/if}
</div>


{#if showAddMemberDialog}
	<AddMemberDialog
		{selectedGroup}
		{groupLabels}
		onClose={handleCloseDialog}
		onMemberAdded={handleMemberAdded}
	/>
{/if}

{#if memberForDetails}
	<MemberDetailDialog
		member={memberForDetails}
		onClose={handleCloseDetailDialog}
		onMemberUpdated={handleMemberUpdated}
	/>
{/if}
