<script lang="ts">
	import { notifications } from '$lib/services/notification.svelte';
</script>

<div class="fixed top-4 right-4 z-5000 ">
	{#each notifications as notification}
		<div class="my-1 px-6 py-3 rounded shadow-lg bg-red-600 text-white font-semibold transition-all">
			{notification.message}
		</div>
	{/each}
</div>

