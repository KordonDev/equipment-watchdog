<script lang="ts">
	import { onMount } from 'svelte';
	import { type Change, getRecentChanges, ChangeAction } from '$lib/services/changes.service';
	import { getOrders, getOrdersOpenOrChanged, type Order } from '$lib/services/order.service';
	import { getMembers, type Member } from '$lib/services/member.service';
	import { translateEquipmentType } from '$lib/services/equipment.service';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';

	let changes: Change[] = [];
	let loading = true;

	let orders: Order[] = [];
	let members: Member[] = [];

	onMount(async () => {
		loading = true;
		[changes, orders, members] = await Promise.all([
			getRecentChanges(),
			getOrdersOpenOrChanged(),
			getMembers()
		]);
		loading = false;
	});

	// Parse "DD/MM/YYYY" to Date
	function parseDDMMYYYY(dateStr: string): Date {
		const [day, month, year] = dateStr.split('/').map(Number);
		return new Date(year, month - 1, day);
	}

	function formatDate(date: Date | string) {
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleDateString();
	}

	function formatTime(date: Date | string) {
		const d = typeof date === 'string' ? new Date(date) : date;
		return d.toLocaleTimeString().slice(0, 5);
	}

	// Helper to get member name by id
	function getMemberName(memberId: number): string {
		const member = members.find(m => m.id === memberId);
		return member ? member.name : `Mitglied #${memberId}`;
	}

	function getEquipmentType(order: Order | undefined, orderId: number): string {
		return order ? translateEquipmentType(order.type) : `#${orderId}`;
	}
	function getOrder(orderId: number): Order | undefined {
		return orders.find(o => o.id === orderId);
	}

	// Group changes by day, newest day first
	$: groupedChanges = (() => {
		const groups: Record<string, Change[]> = {};
		for (const change of changes) {
			const day = formatDate(change.createdAt);
			if (!groups[day]) groups[day] = [];
			groups[day].push(change);
		}
		// Sort days descending, parsing DD/MM/YYYY
		return Object.entries(groups)
			.sort(([a], [b]) => parseDDMMYYYY(b).getTime() - parseDDMMYYYY(a).getTime())
			.map(([day, changes]) => ({
				day,
				changes: changes.sort(
					(a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime()
				)
			}));
	})();

	// Map action to sentence with member name and order type
	function actionToSentence(change: Change): string {
		const memberName = getMemberName(change.memberId);
		const order = getOrder(change.orderId);
		const equipmentType = getEquipmentType(order, change.orderId);

		switch (change.action) {
			case ChangeAction.CreateOrder:
				return `${equipmentType} Größe ${order?.size || '-'} für ${memberName} wurde bestellt.`;
			case ChangeAction.DeleteOrder:
				return `Bestellung ${equipmentType} für ${memberName} wurde gelöscht.`;
			case ChangeAction.UpdateOrder:
				return `Bestellung ${equipmentType} für ${memberName} wurde aktualisiert.`;
			case ChangeAction.OrderToEquipment:
				return `Bestellung ${equipmentType} (${change.equipmentId}) wurde ${memberName} zugeordnet.`;
			case ChangeAction.CreateMember:
				return `${memberName} wurde erstellt.`;
			case ChangeAction.UpdateMember:
				return `${memberName} wurde aktualisiert.`;
			case ChangeAction.DeleteMember:
				return `${memberName} wurde gelöscht.`;
			case ChangeAction.CreateEquipment:
				return `Ausrüstung ${equipmentType} (${change.equipmentId}) wurde erstellt.`;
			case ChangeAction.DeleteEquipment:
				return `Ausrüstung ${change.equipmentId} wurde gelöscht.`;
			case ChangeAction.UpdateEquipmentOnMember:
				if (change.oldEquipmentId) {
					return `${memberName} hat ${equipmentType} ${change.oldEquipmentId} gegen ${change.equipmentId} getauscht.`;
				}
				if (change.equipmentId === 0) {
					return `${memberName} hat ${equipmentType} zurückgegeben.`;
				}
				return `${memberName} hat ${equipmentType} ${change.equipmentId} bekommen.`;
			default:
				return change.action;
		}
	}
</script>

<div class="p-6">
	<BurgerMenu />
	<h1 class="text-2xl font-bold mb-6">Änderungen</h1>
	{#if loading}
		<div>Lade Änderungen...</div>
	{:else if changes.length === 0}
		<div>Keine Änderungen vorhanden.</div>
	{:else}
		{#each groupedChanges as group}
			<div class="mb-8">
				<h2 class="text-lg font-semibold mb-2">{group.day}</h2>
				<ul class="space-y-2">
					{#each group.changes as change}
						<li class="bg-white border rounded px-4 py-2 flex items-center gap-4">
							<span class="text-gray-500 w-16">{formatTime(change.createdAt)}</span>
							<span>{actionToSentence(change)}</span>
						</li>
					{/each}
				</ul>
			</div>
		{/each}
	{/if}
</div>
