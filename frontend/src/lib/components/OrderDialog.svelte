<script lang="ts">
	import { EquipmentType } from '$lib/services/equipment.service';
	import { createOrder, deleteOrder, fulfillOrder, getOrdersForMember, type Order } from '$lib/services/order.service';
	import { getNextGloveId } from '$lib/services/gloveId.service';
	import { onMount } from 'svelte';
	import { showError } from '$lib/services/notification.svelte';

	export let memberName: string;
	export let equipmentLabels: Record<EquipmentType, string>;
	export let onClose: () => void;
	export let memberId: string;
	export let onEquipmentChanged: () => void; // Callback für Equipment-Update

	let orderSizes: Record<EquipmentType, string> = {
		[EquipmentType.Helmet]: '0',
	} as Record<EquipmentType, string>;
	let registrationCodes: Record<EquipmentType, string> = {} as Record<EquipmentType, string>;
	let orders: Order[] = [];
	let loadingOrders = false;
	let nextGloveId: string | null = null;
	let loadingGloveId = false;

	const loadOrders = async () => {
		loadingOrders = true;
		orders = await getOrdersForMember(memberId);
		loadingOrders = false;
	};

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

	onMount(() => {
		loadOrders();
		loadNextGloveId();
	})

	const getOpenOrder = (type: EquipmentType) =>
		orders.find(order => order.type === type && !order.fulfilledAt);

	const handleOrderEquipment = async (e: SubmitEvent,equipmentType: EquipmentType) => {
		e.preventDefault();
		const size = orderSizes[equipmentType]?.trim();
		try {
			const memberIdInternal = parseInt(memberId, 10);
			await createOrder({
				id: 0,
				type: equipmentType,
				createdAt: undefined,
				size,
				memberId: memberIdInternal,
				fulfilledAt: undefined
			});
			if (equipmentType !== EquipmentType.Helmet) {
				orderSizes[equipmentType] = '';
			}
			await loadOrders();
			onEquipmentChanged && onEquipmentChanged(); // Equipment/Order-Update triggern
		} catch (e) {
			showError("Fehler beim Erstellen der Bestellung.");
		}
	};

	const handleDeleteOrder = async (orderId: number) => {
		await deleteOrder(orderId);
		await loadOrders();
		onEquipmentChanged && onEquipmentChanged(); // Equipment/Order-Update triggern
	};

	const handleFulfillOrder = async (e: SubmitEvent,order: Order, equipmentType: EquipmentType) => {
		e.preventDefault()
		const regCode = registrationCodes[equipmentType]?.trim();
		if (!regCode) {
			alert('Bitte eine Ausrüstungsnummer angeben.');
			return;
		}

		await fulfillOrder(order, regCode);
		registrationCodes[equipmentType] = '';
		await loadOrders();
		onEquipmentChanged && onEquipmentChanged(); // Equipment-Update triggern
	};

	function formatDate(date: Date | undefined) {
		if (!date) return '';
		return date.toLocaleDateString();
	}
</script>

	<div class="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-60">
		<div class="bg-white rounded-lg p-6 w-full max-w-md shadow-lg relative">
			<div class="flex items-center justify-between mb-4">
				<h3 class="text-lg font-bold">Bestellungen {memberName}</h3>
				<button
					type="button"
					on:click={onClose}
					class="text-gray-400 hover:text-gray-600 text-2xl leading-none"
				>
					×
				</button>
			</div>
			{#if loadingOrders}
				<div class="text-center py-4">Lade Bestellungen...</div>
			{:else}
				<div class="space-y-4">
					{#each Object.values(EquipmentType) as equipmentType}
						{@const openOrder = getOpenOrder(equipmentType)}
						<div class="flex flex-col gap-1">
							<div class="flex items-center gap-2">
								<span class="w-28">{equipmentLabels[equipmentType]}</span>
								{#if openOrder}
									{#if equipmentType !== EquipmentType.Helmet}
										<span class="text-xs text-gray-600 bg-gray-100 rounded px-2 py-0.5" title="Bestellte Größe">
											Größe: {openOrder.size}
										</span>
									{/if}
									<span class="text-xs text-gray-500 bg-gray-50 rounded px-2 py-0.5" title="Bestelldatum">
										{formatDate(openOrder.createdAt)}
									</span>
									<div class="flex-1"></div>
									<button
										type="button"
										class="px-3 py-1 bg-red-600 text-white rounded hover:bg-red-700 transition-colors ml-auto"
										on:click={() => handleDeleteOrder(openOrder.id)}
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
								<form onsubmit={e => handleFulfillOrder(e, openOrder, equipmentType)} class="flex items-center gap-2 mt-1">
									<input
										type="text"
										placeholder="Ausrüstungsnummer"
										required
										class="flex-1 px-2 py-1 border border-gray-300 rounded focus:border-blue-500"
										bind:value={registrationCodes[equipmentType]}
									/>
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
								{#if equipmentType === EquipmentType.Gloves && loadingGloveId}
									<div class="text-xs text-blue-600 mt-1">
										⏳ Lade nächste Handschuh-ID...
									</div>
								{/if}
							{/if}
						</div>
					{/each}
				</div>
			{/if}
			<div class="flex justify-end mt-6">
				<button
					type="button"
					on:click={onClose}
					class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
				>
					Schließen
				</button>
			</div>
		</div>
	</div>
