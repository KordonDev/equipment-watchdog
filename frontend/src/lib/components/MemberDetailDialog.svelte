<script lang="ts">
	import {
		type Member,
		updateMember,
		deleteMember,
		Group,
		getMember,
		groupLabels,
		equipmentForGroup
	} from '$lib/services/member.service';
	import {
		saveEquipmentForMember,
		type Equipment,
		equipmentLabels,
		EquipmentType, randomRegistrationCode
	} from '$lib/services/equipment.service';
	import { getOrdersForMember, type Order } from '$lib/services/order.service';
	import { getNextGloveId } from '$lib/services/gloveId.service';
	import OrderDialog from './OrderDialog.svelte';
	import { showError } from '$lib/services/notification.svelte';

	interface Props {
		member: Member | null;
		onClose: (memberId?: string) => void;
		onMemberUpdated: (member: Member) => void;
	}

	let { member, onClose, onMemberUpdated }: Props = $props();

	let editingMember: Member | null = $state(JSON.parse(JSON.stringify(member)));// Deep copy
	let saving = $state(false);
	let tempRegistrationCodes: Record<EquipmentType, string> = $state({} as Record<EquipmentType, string>);
	let nextGloveId: string | null = $state(null);
	let loadingGloveId = $state(false);

	Object.values(EquipmentType).forEach(type => {
		tempRegistrationCodes[type] = editingMember?.equipments[type]?.registrationCode || '';
	});
	if (tempRegistrationCodes[EquipmentType.Helmet] === '') {
		tempRegistrationCodes[EquipmentType.Helmet] = randomRegistrationCode();
	}

	$effect(() => {
		if (editingMember) {
			loadNextGloveId();
		}
	});

	const loadNextGloveId = async () => {
		try {
			loadingGloveId = true;
			nextGloveId = (await getNextGloveId()).nextId;
		} catch (error) {
			console.error('Failed to load next glove ID:', error);
			nextGloveId = null;
		} finally {
			loadingGloveId = false;
		}
	};

	const handleClose = () => {
		editingMember = null;
		onClose();
	};

	const addEquipment = async (equipmentType: EquipmentType) => {
		if (!editingMember) return;

		try {
			// Use new single equipment save endpoint
			const newEquipment = await saveEquipmentForMember(
				editingMember.id,
				equipmentType,
				{
					id: 0,
					type: equipmentType,
					registrationCode: tempRegistrationCodes[equipmentType] || '',
					size: '42',
					memberId: editingMember.id
				}
			);

			// Update local state
			editingMember.equipments[equipmentType] = newEquipment;

			// Refresh member data from server
			const updatedMember = await getMember(editingMember.id.toString());
			editingMember = updatedMember;
			onMemberUpdated(updatedMember);
		} catch (error) {
			showError('Fehler beim Hinzufügen der Ausrüstung.');
			console.error(error);
		}
	};

	const removeEquipment = async (equipmentType: EquipmentType) => {
		if (!editingMember) return;

		try {
			// Remove from local state (backend handles on updateMember)
			delete editingMember.equipments[equipmentType];

			// Update member to unassign equipment
			const updatedMember = await updateMember(editingMember);
			tempRegistrationCodes[equipmentType] = '';
			if (equipmentType === EquipmentType.Helmet) {
				tempRegistrationCodes[equipmentType] = randomRegistrationCode();
			}
			editingMember = updatedMember;
			onMemberUpdated(updatedMember);
		} catch (error) {
			showError('Fehler beim Entfernen der Ausrüstung.');
			console.error(error);
		}
	};

	const updateEquipmentNumber = (equipmentType: EquipmentType, number: string) => {
		if (!editingMember) return;

		// Only allow updating registration codes for equipment that hasn't been saved yet
		if (!editingMember.equipments[equipmentType]) {
			tempRegistrationCodes[equipmentType] = number;
		}
	};

	let groupChanged = $state(false);

	const handleGroupInput = (event: Event) => {
		if (!editingMember) return;
		const select = event.target as HTMLSelectElement;
		const newGroup = select.value as Group;
		groupChanged = editingMember.group !== newGroup;
		editingMember.group = newGroup;
	};

	const handleGroupSave = async (e: SubmitEvent) => {
		e.preventDefault();
		if (!editingMember) return;
		try {
			const updatedMember = await updateMember(editingMember);
			editingMember = updatedMember;
			if (member) {
				member.group = updatedMember.group;
			}
			onMemberUpdated(updatedMember);
			groupChanged = false;
		} catch (error) {
			console.error(error);
			showError('Fehler beim Speichern der Gruppe.');
		}
	};

	const handleDeleteMember = async () => {
		if (!editingMember) return;
		if (!confirm(`Soll ${editingMember.name} wirklich gelöscht werden?`)) return;
		try {
			await deleteMember(editingMember.id);
			onClose(editingMember.id); // Pass deleted member ID
			// Optionally, you can notify parent to refresh the member list here
		} catch (error) {
			console.error(error);
			showError("Fehler beim Löschen des Mitglieds.");
		}
	};

	let showOrderDialog = $state(false);

	const handleOpenOrderDialog = () => {
		showOrderDialog = true;
	};

	const handleCloseOrderDialog = () => {
		showOrderDialog = false;
	};

	let orders: Order[] = $state([]);

	$effect(() => {
		if (editingMember !== null) {
			getOrdersForMember(editingMember.id.toString())
				.then(o => orders = o)
		}
	})

	const hasOrderForType = (type: EquipmentType) =>
		orders.some(order => order.type === type && !order.fulfilledAt);

	const handleEquipmentChanged = async () => {
		if (!editingMember) return;
		const updated = await getMember(editingMember.id.toString());
		editingMember = updated;
		onMemberUpdated(updated);
	};
</script>

{#if editingMember}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[80vh] overflow-y-auto">
			<div class="flex items-start justify-between mb-4">
				<h2 class="text-xl font-bold">{editingMember.name} bearbeiten</h2>
				<button
					type="button"
					onclick={handleClose}
					class="text-gray-400 hover:text-gray-600 text-2xl leading-none"
				>
					×
				</button>
			</div>

			<!-- Order Equipment Button -->
			<div class="flex justify-end mb-4">
				<button
					type="button"
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
					onclick={handleOpenOrderDialog}
				>
					Bestellungen
				</button>
			</div>

			<div class="space-y-4 mb-4">
				{#each equipmentForGroup[member?.group] as equipmentType}
					{@const equipment = editingMember.equipments[equipmentType]}
					<form class="border rounded-lg p-4" onsubmit={() => equipment ? removeEquipment(equipmentType) : addEquipment(equipmentType)}>
						<div class="flex items-center justify-between mb-2">
							<div class="text-sm font-medium text-gray-700 flex items-center gap-2">
								{equipmentLabels[equipmentType]}
								{#if hasOrderForType(equipmentType)}
									<span class="ml-2 px-2 py-0.5 text-xs rounded bg-yellow-100 text-yellow-800" title="Bestellt">🛒 Bestellt</span>
								{/if}
							</div>
							<button
								type="submit"
								class="px-3 py-1 text-xs rounded {equipment ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-green-100 text-green-700 hover:bg-green-200'}"
							>
								{equipment ? 'Entfernen' : 'Hinzufügen'}
							</button>
						</div>

						<div>
							{#if equipmentType !== EquipmentType.Helmet}
								<label class="block text-xs text-gray-500 mb-1">Ausrüstungsnummer
								<input
									type="text"
									required
									value={equipment ? equipment.registrationCode : tempRegistrationCodes[equipmentType] || ''}
									oninput={(e) => updateEquipmentNumber(equipmentType, e.target?.value || '')}
									disabled={!!equipment}
									class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50 {equipment ? 'bg-gray-100 text-gray-500 cursor-not-allowed' : ''}"
								/>
								</label>
							{/if}
							{#if equipmentType === EquipmentType.Gloves && !equipment && nextGloveId}
								<div class="text-xs text-green-600 mt-1">
									✓ Nächste verfügbare Handschuh-ID: {nextGloveId}
								</div>
							{/if}
							{#if equipmentType === EquipmentType.Gloves && !equipment && loadingGloveId}
								<div class="text-xs text-blue-600 mt-1">
									⏳ Lade nächste Handschuh-ID...
								</div>
							{/if}
						</div>
					</form>
				{/each}
			</div>

			<!-- Group Selector with Save Button -->
			<form onsubmit={handleGroupSave} class="mb-4 flex items-end gap-2">
				<label class="block text-xs text-gray-500 mb-1 flex-grow">Gruppe
					<select
						class="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
						bind:value={editingMember.group}
						oninput={handleGroupInput}
					>
						{#each Object.values(Group) as group}
							<option value={group}>{groupLabels[group]}</option>
						{/each}
					</select>
				</label>
				<button
					type="submit"
					style="margin-bottom: 4px; padding-bottom: 3px; padding-top: 3px;"
					class="px-3 py-1 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors"
					disabled={!groupChanged}
				>
					Speichern
				</button>
			</form>
			
			<div class="flex justify-end mt-2">
				<button
					type="button"
					onclick={handleDeleteMember}
					class="px-4 py-2 text-white bg-red-600 rounded hover:bg-red-700 transition-colors"
					disabled={saving}
				>
					Mitglied löschen
				</button>
			</div>
			<div class="flex justify-end space-x-2 mt-6">
				<button
					type="button"
					onclick={handleClose}
					class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
					disabled={saving}
				>
					Schließen
				</button>
			</div>
		</div>
	</div>

	{#if showOrderDialog}
	<OrderDialog
		member={editingMember}
		equipmentLabels={equipmentLabels}
		onClose={handleCloseOrderDialog}
		onEquipmentChanged={handleEquipmentChanged}
	/>
	{/if}
{/if}
