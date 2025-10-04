<script lang="ts">
	import { createMember, Group, type Member } from '$lib/services/member.service';
	import { showError } from '$lib/services/notification.svelte';

	interface Props {
		selectedGroup: Group;
		groupLabels: Record<Group, string>;
		onClose: () => void;
		onMemberAdded: (member: Member) => void;
	}

	let { selectedGroup, groupLabels, onClose, onMemberAdded }: Props = $props();

	let newMemberName = $state('');
	let addingMember = $state(false);

	const handleSubmitMember = async (e: SubmitEvent) => {
		e.preventDefault();
		if (!newMemberName.trim()) return;

		addingMember = true;
		try {
			const newMember: Member = {
				id: 0, // Will be set by backend
				name: newMemberName.trim(),
				group: selectedGroup,
				equipments: {}
			};

			const createdMember = await createMember(newMember);
			onMemberAdded(createdMember);
			newMemberName = '';
		} catch (error) {
			console.error(error);
			showError('Fehler beim Anlegen des Mitglieds.');
		} finally {
			addingMember = false;
		}
	};

	const handleCancel = () => {
		newMemberName = '';
		onClose();
	};
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
	<div class="bg-white rounded-lg p-6 w-full max-w-md">
		<div class="flex items-start justify-between mb-4">
			<h2 class="text-xl font-bold mb-4">Neues Mitglied hinzufügen</h2>
			<button
				type="button"
				onclick={handleCancel}
				class="text-gray-400 hover:text-gray-600 text-2xl leading-none"
			>
				×
			</button>
		</div>
		<form onsubmit={handleSubmitMember}>

			<p class="text-sm text-gray-600 mb-4">Gruppe: {groupLabels[selectedGroup]}</p>

			<div class="mb-4">
				<label for="memberName" class="block mb-2 text-sm font-medium text-gray-700">Name:</label>
				<input
					id="memberName"
					type="text"
					bind:value={newMemberName}
					class="block w-full rounded border border-gray-300 px-3 py-2 focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
					placeholder="Name eingeben..."
				/>
			</div>

			<div class="flex justify-end space-x-2">
				<button
					type="button"
					onclick={handleCancel}
					class="px-4 py-2 text-gray-600 bg-gray-200 rounded hover:bg-gray-300 transition-colors"
					disabled={addingMember}
				>
					Abbrechen
				</button>
				<button
					type="submit"
					class="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors disabled:opacity-50"
					disabled={addingMember || !newMemberName.trim()}
				>
					{addingMember ? 'Hinzufügen...' : 'Hinzufügen'}
				</button>
			</div>
		</form>
	</div>
</div>
