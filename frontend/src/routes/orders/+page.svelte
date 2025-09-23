<script lang="ts">
	import { onMount } from 'svelte';
	import { getOrders, type Order } from '$lib/services/order.service';
	import { EquipmentType } from '$lib/services/equipment.service';
	import { getMembers, type Member } from '$lib/services/member.service';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';

	let orders: Order[] = [];
	let members: Member[] = [];
	let loading = true;

	const equipmentLabels: Record<EquipmentType, string> = {
		[EquipmentType.Helmet]: 'Helm',
		[EquipmentType.Jacket]: 'Jacke',
		[EquipmentType.Gloves]: 'Handschuhe',
		[EquipmentType.Trousers]: 'Hose',
		[EquipmentType.Boots]: 'Stiefel',
		[EquipmentType.TShirt]: 'T-Shirt'
	};

	onMount(async () => {
		loading = true;
		orders = await getOrders();
		members = await getMembers();
		loading = false;
	});

	function getMemberName(memberId: number) {
		const member = members.find(m => m.id === memberId);
		return member ? member.name : 'Unbekannt';
	}

	function formatDate(date: Date | string | undefined) {
		if (!date) return '';
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString();
	}

	// Create segments for each equipment type
	const equipmentTypes = Object.values(EquipmentType);

	function getOrdersByType(type: EquipmentType) {
		return orders.filter(order => order.type === type);
	}
</script>

<div class="p-6">
	<BurgerMenu />
	<h1 class="text-2xl font-bold mb-6">Bestellungen</h1>
	{#if loading}
		<div>Lade Bestellungen...</div>
	{:else}
		<!-- Fast scroll links for each equipment type -->
		<div class="mb-4 flex flex-wrap gap-4">
			{#each equipmentTypes as type}
				<a href={"#type-" + type} class="px-3 py-1 bg-blue-100 text-blue-800 rounded hover:bg-blue-200 transition-colors">
					{equipmentLabels[type]}
				</a>
			{/each}
		</div>
		<!-- Segments for each equipment type -->
		{#each equipmentTypes as type}
			<div id={"type-" + type} class="mb-8">
				<h2 class="text-xl font-semibold mb-2">{equipmentLabels[type]}</h2>
				<table class="min-w-full bg-white border border-gray-200 rounded">
					<thead>
						<tr>
							<th class="px-4 py-2 border-b">Größe</th>
							<th class="px-4 py-2 border-b">Bestelldatum</th>
							<th class="px-4 py-2 border-b">Für</th>
						</tr>
					</thead>
					<tbody>
						{#each getOrdersByType(type) as order}
							<tr>
								<td class="px-4 py-2 border-b">{order.size}</td>
								<td class="px-4 py-2 border-b">{formatDate(order.createdAt)}</td>
								<td class="px-4 py-2 border-b">{getMemberName(order.memberId)}</td>
							</tr>
						{/each}
						{#if getOrdersByType(type).length === 0}
							<tr>
								<td colspan="3" class="px-4 py-2 border-b text-gray-400">Keine Bestellungen vorhanden.</td>
							</tr>
						{/if}
					</tbody>
				</table>
			</div>
		{/each}
	{/if}
</div>
