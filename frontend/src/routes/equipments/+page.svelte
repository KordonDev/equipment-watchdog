<script lang="ts">
	import { onMount } from 'svelte';
	import { EquipmentType, type Equipment, getAllEquipment } from '$lib/services/equipment.service';
	import { getMembers, type Member } from '$lib/services/member.service';
	import BurgerMenu from '$lib/components/BurgerMenu.svelte';

	let equipment: Equipment[] = [];
	let members: Member[] = [];
	let loading = true;

	const equipmentLabels: Record<EquipmentType, string> = {
		[EquipmentType.Gloves]: 'Handschuhe',
		[EquipmentType.Jacket]: 'Jacke',
		[EquipmentType.Trousers]: 'Hose',
		[EquipmentType.Boots]: 'Stiefel',
		[EquipmentType.TShirt]: 'TShirt',
		[EquipmentType.Helmet]: 'Helm'
	};

	onMount(async () => {
		loading = true;
		equipment = await getAllEquipment();
		members = await getMembers();
		loading = false;
	});

	function getOwner(e: Equipment) {
		const member = members.find(m => m.id === e.memberId);
		return member ? member.name : 'Unbekannt';
	}

	function sortedEquipmentByType(type: EquipmentType) {
		return equipment.filter(e => e.type === type);
	}
</script>

<div class="p-6">
	<BurgerMenu />
	<h1 class="text-2xl font-bold mb-6">Ausrüstungsliste</h1>
	{#if loading}
		<div>Lade Ausrüstung...</div>
	{:else}
		<!-- Fast scroll links -->
		<div class="mb-4 flex flex-wrap gap-4">
			{#each Object.values(EquipmentType) as type}
				<a href={"#type-" + type} class="px-3 py-1 bg-blue-100 text-blue-800 rounded hover:bg-blue-200 transition-colors">
					{equipmentLabels[type]}
				</a>
			{/each}
		</div>
		<!-- Equipment lists by type -->
		{#each Object.values(EquipmentType) as type}
			<div id={"type-" + type} class="mb-8">
				<h2 class="text-xl font-semibold mb-2">{equipmentLabels[type]}</h2>
				<table class="min-w-full bg-white border border-gray-200 rounded">
					<thead>
						<tr>
							<th class="px-4 py-2 border-b">Code</th>
							<th class="px-4 py-2 border-b">Größe</th>
							<th class="px-4 py-2 border-b">Besitzer</th>
						</tr>
					</thead>
					<tbody>
						{#each sortedEquipmentByType(type) as e}
							<tr>
								<td class="px-4 py-2 border-b">{e.registrationCode}</td>
								<td class="px-4 py-2 border-b">{e.size}</td>
								<td class="px-4 py-2 border-b">{getOwner(e)}</td>
							</tr>
						{/each}
						{#if sortedEquipmentByType(type).length === 0}
							<tr>
								<td colspan="3" class="px-4 py-2 border-b text-gray-400">Keine Ausrüstung vorhanden.</td>
							</tr>
						{/if}
					</tbody>
				</table>
			</div>
		{/each}
	{/if}
</div>
