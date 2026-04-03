<script lang="ts">
	import {
		type Member,
		updateMember,
		deleteMember,
		Group,
		groupLabels,
		equipmentForGroup,
		saveEquipmentForMember,
		removeEquipmentFromMember
	} from '$lib/services/member.service';
	import {
		type Equipment,
		equipmentLabels,
		EquipmentType,
		randomRegistrationCode
	} from '$lib/services/equipment.service';
	import { getOrdersForMember, type Order } from '$lib/services/order.service';
	import { getNextGloveId } from '$lib/services/gloveId.service';
	import { showError } from '$lib/services/notification.svelte';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let member: Member = $state(JSON.parse(JSON.stringify(data.member)));

	let saving = $state(false);
	let tempRegistrationCodes: Record<EquipmentType, string> = $state({} as Record<EquipmentType, string>);
	let tempGroup = $state(data.member.group);
	let groupChanged = $state(false);
	let nextGloveId: string | null = $state(null);
	let loadingGloveId = $state(false);
	let orders: Order[] = $state([]);

	Object.values(EquipmentType).forEach(type => {
		tempRegistrationCodes[type] = member.equipments[type]?.registrationCode || '';
	});
	if (tempRegistrationCodes[EquipmentType.Helmet] === '') {
		tempRegistrationCodes[EquipmentType.Helmet] = randomRegistrationCode();
	}

	onMount(() => {
		getOrdersForMember(member.id.toString()).then(o => (orders = o));
		loadNextGloveId();
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

	const addEquipment = async (equipmentType: EquipmentType) => {
		const registrationCode = tempRegistrationCodes[equipmentType];
		if (!registrationCode) return;

		try {
			const newEquipment = await saveEquipmentForMember(member.id, equipmentType, {
				id: 0,
				type: equipmentType,
				registrationCode,
				size: '42',
				memberId: member.id
			});
			member.equipments[equipmentType] = newEquipment;
			member = { ...member };
		} catch (error) {
			showError('Fehler beim Hinzufügen der Ausrüstung.');
			console.error(error);
		}
	};

	const removeEquipment = async (equipmentType: EquipmentType) => {
		try {
			await removeEquipmentFromMember(member.id, equipmentType);
			delete member.equipments[equipmentType];
			member = { ...member };
			tempRegistrationCodes[equipmentType] = '';
			if (equipmentType === EquipmentType.Helmet) {
				tempRegistrationCodes[EquipmentType.Helmet] = randomRegistrationCode();
			}
		} catch (error) {
			showError('Fehler beim Entfernen der Ausrüstung.');
			console.error(error);
		}
	};

	const handleEquipmentNumberInput = (equipmentType: EquipmentType, number: string) => {
		if (!member.equipments[equipmentType]) {
			tempRegistrationCodes[equipmentType] = number;
		}
	};

	const handleGroupInput = (event: Event) => {
		const select = event.target as HTMLSelectElement;
		const newGroup = select.value as Group;
		groupChanged = member.group !== newGroup;
	};

	const handleGroupSave = async (e: SubmitEvent) => {
		e.preventDefault();
		try {
			const updatedMember = await updateMember({ ...member, group: tempGroup });
			member.group = updatedMember.group;
			member = { ...member };
			groupChanged = false;
		} catch (error) {
			console.error(error);
			showError('Fehler beim Speichern der Gruppe.');
		}
	};

	const handleDeleteMember = async () => {
		if (!confirm(`Soll ${member.name} wirklich gelöscht werden?`)) return;
		try {
			await deleteMember(member.id.toString());
			goto('/members');
		} catch (error) {
			console.error(error);
			showError('Fehler beim Löschen des Mitglieds.');
		}
	};

	const hasOrderForType = (type: EquipmentType) =>
		orders.some(order => order.type === type && !order.fulfilledAt);

	const hasAssignedEquipment = () =>
		Object.values(member.equipments).some(e => !!e);
</script>

<div class="p-6">
	<BurgerMenu />
	<div class="max-w-2xl mx-auto">
		<div class="flex items-center gap-4 mb-6">
			<button
				type="button"
				onclick={() => goto('/members')}
				class="text-gray-500 hover:text-gray-700 text-sm flex items-center gap-1"
			>
				← Zurück
			</button>
			<h1 class="text-2xl font-bold">{member.name} bearbeiten</h1>
		</div>

		<!-- Orders Button -->
		<div class="flex justify-end mb-4">
			<button
				type="button"
				class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
				onclick={() => goto(`/members/${member.id}/orders`)}
			>
				Bestellungen
			</button>
		</div>

		<div class="space-y-4 mb-4">
			{#each equipmentForGroup[member.group] as equipmentType}
				{@const equipment = member.equipments[equipmentType]}
				<form
					class="border rounded-lg p-4"
					onsubmit={() => (equipment ? removeEquipment(equipmentType) : addEquipment(equipmentType))}
				>
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
									value={tempRegistrationCodes[equipmentType] || ''}
									oninput={(e) => handleEquipmentNumberInput(equipmentType, (e.target as HTMLInputElement)?.value || '')}
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
					bind:value={tempGroup}
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

		{#if member.group === Group.GONE}
			<div class="flex justify-end mt-2">
				<button
					type="button"
					onclick={handleDeleteMember}
					class="px-4 py-2 text-white bg-red-600 rounded hover:bg-red-700 transition-colors"
					disabled={saving || hasAssignedEquipment()}
				>
					Mitglied löschen
				</button>
			</div>
		{/if}

		<div class="flex justify-end space-x-2 mt-6">
			<button
				type="button"
				onclick={() => goto('/members')}
				class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
				disabled={saving}
			>
				Schließen
			</button>
		</div>
	</div>
</div>
