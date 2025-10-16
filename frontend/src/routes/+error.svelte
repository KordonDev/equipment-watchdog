<script lang="ts">
	import { page } from '$app/state';
	import { goto } from '$app/navigation';

	if (page && page.status === 404) {
		goto("/members", { replaceState: true });
	}
</script>
