<script lang="ts">
	import {
		type Equipment,
		equipmentLabels,
		EquipmentType,
		randomRegistrationCode
	} from '$lib/services/equipment.service';
	import {
		createOrder,
		deleteOrder,
		fulfillOrder,
		type Order
	} from '$lib/services/order.service';
	import { equipmentForGroup, saveEquipmentForMember } from '$lib/services/member.service';
	import { showError } from '$lib/services/notification.svelte';
	import { formatToDate } from '$lib/services/timeHelper';
	import { goto } from '$app/navigation';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let member = $state(data.member);
	let orders: Order[] = $state(data.orders);
	let nextGloveId: string | null = $state(data.nextGloveId);

	let orderSizes: Record<EquipmentType, string> = $state({
		[EquipmentType.Helmet]: '0'
	}) as Record<EquipmentType, string>;

	let registrationCodes: Record<EquipmentType, string> = $state({
		[EquipmentType.Helmet]: randomRegistrationCode()
	}) as Record<EquipmentType, string>;

	const getOpenOrder = (type: EquipmentType) =>
		orders.find(order => order.type === type && !order.fulfilledAt);

	const handleOrderEquipment = async (e: SubmitEvent, equipmentType: EquipmentType) => {
		e.preventDefault();
		const size = orderSizes[equipmentType]?.trim();
		try {
			const newOrder = await createOrder({
				id: 0,
				type: equipmentType,
				createdAt: undefined,
				size,
				memberId: parseInt(member.id, 10),
				fulfilledAt: undefined
			});
			if (equipmentType !== EquipmentType.Helmet) {
				orderSizes[equipmentType] = '';
			}
			orders = [...orders, newOrder];
		} catch (err) {
			console.error(err);
			showError('Fehler beim Erstellen der Bestellung.');
		}
	};

	const handleDeleteOrder = async (orderId: number) => {
		await deleteOrder(orderId);
		orders = orders.filter(order => order.id !== orderId);
	};

	const handleFulfillOrder = async (e: SubmitEvent, order: Order, equipmentType: EquipmentType) => {
		e.preventDefault();
		const regCode = registrationCodes[equipmentType]?.trim();
		if (!regCode) {
			alert('Bitte eine Ausrüstungsnummer angeben.');
			return;
		}
		try {
			const equipment: Equipment = await fulfillOrder(order, regCode);
			registrationCodes[equipmentType] = '';
			orders = orders.filter(o => o.id !== order.id);
			member = {
				...member,
				equipments: { ...member.equipments, [equipment.type]: equipment }
			};
		} catch (error) {
			console.error(error);
			showError('Fehler beim Erfüllen der Bestellung.');
		}
	};
</script>

<div class="p-6">
	<BurgerMenu />
	<div class="max-w-md mx-auto">
		<div class="flex items-center gap-4 mb-6">
			<button
				type="button"
				onclick={() => goto(`/members/${member.id}`)}
				class="text-gray-500 hover:text-gray-700 text-sm flex items-center gap-1"
			>
				← Zurück
			</button>
			<h1 class="text-2xl font-bold">Bestellungen {member.name}</h1>
		</div>

		<div class="space-y-4">
			{#each equipmentForGroup[member.group] as equipmentType}
				{@const openOrder = getOpenOrder(equipmentType)}
				<div class="border rounded-lg p-4 flex flex-col gap-1">
					<div class="flex items-center gap-2">
						<span class="w-28">{equipmentLabels[equipmentType]}</span>
						{#if openOrder}
							{#if equipmentType !== EquipmentType.Helmet}
								<span class="text-xs text-gray-600 bg-gray-100 rounded px-2 py-0.5" title="Bestellte Größe">
									Größe: {openOrder.size}
								</span>
							{/if}
							<span class="text-xs text-gray-500 bg-gray-50 rounded px-2 py-0.5" title="Bestelldatum">
								{formatToDate(openOrder.createdAt)}
							</span>
							<div class="flex-1"></div>
							<button
								type="button"
								class="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700 transition-colors ml-auto"
								onclick={() => handleDeleteOrder(openOrder.id)}
							>
								Löschen
							</button>
						{:else}
							<form onsubmit={e => handleOrderEquipment(e, equipmentType)}>
								{#if equipmentType !== EquipmentType.Helmet}
									<input
										type="text"
										placeholder="Größe"
										class="flex-1 px-2 py-1 border border-gray-300 rounded focus:border-blue-500"
										required
										bind:value={orderSizes[equipmentType]}
									/>
								{:else}
									<div class="flex-1"></div>
								{/if}
								<button
									type="submit"
									class="px-3 py-1 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
								>
									Bestellen
								</button>
							</form>
						{/if}
					</div>

					{#if openOrder}
						<form
							onsubmit={e => handleFulfillOrder(e, openOrder, equipmentType)}
							class="flex items-center gap-2 mt-1"
						>
							{#if equipmentType !== EquipmentType.Helmet}
								<input
									type="text"
									placeholder="Ausrüstungsnummer"
									required
									class="flex-1 px-2 py-1 border border-gray-300 rounded focus:border-blue-500"
									bind:value={registrationCodes[equipmentType]}
								/>
							{:else}
								<div class="flex-1"></div>
							{/if}
							<button
								type="submit"
								class="px-3 py-1 bg-green-600 text-white rounded hover:bg-green-700 transition-colors"
							>
								Erfüllen
							</button>
						</form>
						{#if equipmentType === EquipmentType.Gloves && nextGloveId}
							<div class="text-xs text-green-600 mt-1">
								✓ Nächste verfügbare Handschuh-ID: {nextGloveId}
							</div>
						{/if}
					{/if}
				</div>
			{/each}
		</div>

		<div class="flex justify-end mt-6">
			<button
				type="button"
				onclick={() => goto(`/members/${member.id}`)}
				class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
			>
				Schließen
			</button>
		</div>
	</div>
</div>
