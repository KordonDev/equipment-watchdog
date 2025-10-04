<script lang="ts">
	import { notifications } from '$lib/services/notification.svelte';
</script>

<div class="fixed top-4 right-4 z-5000 ">
	{#each notifications as notification}
		<div class={`my-1 px-6 py-3 rounded shadow-lg text-white font-semibold transition-all ${notification.type}`}>
			{notification.message}
		</div>
	{/each}
</div>

<style>
	.success {
		background-color: #4caf50;
	}
	.error {
		background-color: #f44336;
	}
	.info {
		background-color: #2196f3;
	}
</style>
