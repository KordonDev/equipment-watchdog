<script lang="ts">
	import { type Member, updateMember } from '$lib/services/member.service';
	import { createEquipment, deleteEquipment, type Equipment, EquipmentType } from '$lib/services/equipment.service';

	interface Props {
		member: Member | null;
		onClose: () => void;
		onMemberUpdated: (member: Member) => void;
	}

	let { member, onClose, onMemberUpdated }: Props = $props();

	let editingMember: Member | null = $state(JSON.parse(JSON.stringify(member)));// Deep copy
	let saving = $state(false);
	let tempRegistrationCodes: Record<EquipmentType, string> = $state({} as Record<EquipmentType, string>);
	Object.values(EquipmentType).forEach(type => {
		tempRegistrationCodes[type] = editingMember?.equipments[type]?.registrationCode || '';
	});

	const equipmentLabels = {
		[EquipmentType.Helmet]: 'Helm',
		[EquipmentType.Jacket]: 'Jacke',
		[EquipmentType.Gloves]: 'Handschuhe',
		[EquipmentType.Trousers]: 'Hose',
		[EquipmentType.Boots]: 'Stiefel',
		[EquipmentType.TShirt]: 'T-Shirt'
	};

	const handleClose = () => {
		editingMember = null;
		onClose();
	};

	const addEquipment = async (equipmentType: EquipmentType) => {
		if (!editingMember) return;

		try {
			// Add equipment with registration code from input
			const newEquipment: Equipment = {
				id: 0,
				type: equipmentType,
				registrationCode: tempRegistrationCodes[equipmentType] || '',
				size: '42'
			};
			editingMember.equipments[equipmentType] = await createEquipment(newEquipment);

			// Save member after equipment change
			const updatedMember = await updateMember(editingMember);
			editingMember = updatedMember;
			onMemberUpdated(updatedMember);
		} catch (error) {
			console.error('Failed to add equipment or update member:', error);
			return;
		}
		editingMember = { ...editingMember }; // Trigger reactivity
	};

	const removeEquipment = async (equipmentType: EquipmentType) => {
		if (!editingMember) return;

		try {
			// Delete equipment via backend call
			const equipmentToDelete = editingMember.equipments[equipmentType];
			if (equipmentToDelete) {
				await deleteEquipment(equipmentToDelete.id);
			}
			// Remove equipment from member
			delete editingMember.equipments[equipmentType];

			// Save member after equipment change
			const updatedMember = await updateMember(editingMember);
			editingMember = updatedMember;
			onMemberUpdated(updatedMember);
		} catch (error) {
			console.error('Failed to remove equipment or update member:', error);
			return;
		}
		editingMember = { ...editingMember }; // Trigger reactivity
	};

	const updateEquipmentNumber = (equipmentType: EquipmentType, number: string) => {
		if (!editingMember) return;

		// Only allow updating registration codes for equipment that hasn't been saved yet
		if (!editingMember.equipments[equipmentType]) {
			tempRegistrationCodes[equipmentType] = number;
		}
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

			<div class="space-y-4">
				{#each Object.values(EquipmentType) as equipmentType}
					{@const equipment = editingMember.equipments[equipmentType]}
					<form class="border rounded-lg p-4" onsubmit={() => equipment ? removeEquipment(equipmentType) : addEquipment(equipmentType)}>
						<div class="flex items-center justify-between mb-2">
							<div class="text-sm font-medium text-gray-700">
								{equipmentLabels[equipmentType]}
							</div>
							<button
								type="submit"
								class="px-3 py-1 text-xs rounded {equipment ? 'bg-red-100 text-red-700 hover:bg-red-200' : 'bg-green-100 text-green-700 hover:bg-green-200'}"
							>
								{equipment ? 'Entfernen' : 'Hinzufügen'}
							</button>
						</div>

						<div>
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
						</div>
					</form>
				{/each}
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
{/if}
