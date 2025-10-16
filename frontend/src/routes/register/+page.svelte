<script lang="ts">
  import { register } from "$lib/services/authentication";
	import { goto } from '$app/navigation';
	import { routes } from '$lib/routes';
	import { showError } from '$lib/services/notification.svelte';

  let username = '';
  let loading = false;
  let error: string | null = null;

  const handleRegister = async (event: SubmitEvent) => {
    event.preventDefault();
    loading = true;
    error = null;
    try {
      await register(username);
			goto(`${routes.login}?message=Registrierung erfolgreich! Bitte einloggen.&username=${encodeURIComponent(username)}`);
      // Optionally redirect or show success message here
    } catch (e) {
      showError("Registrierung fehlgeschlagen. Bitte versuche es erneut.")
    } finally {
      loading = false;
    }
  };
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
  <div class="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
    <h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Registrieren</h2>
    <form on:submit={handleRegister} class="space-y-4">
      <div>
        <label for="username" class="block mb-2 text-sm font-medium text-gray-700">Benutzername:</label>
        <input
          required
          class="block w-full rounded border border-gray-300 px-3 py-2 focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
          id="username"
          bind:value={username}
          type="text"
        />
      </div>
      <button
        type="submit"
        class="w-full rounded bg-purple-600 px-4 py-2 text-white font-semibold hover:bg-purple-700 transition-colors duration-150 disabled:opacity-50"
        disabled={loading}
      >
        Registrieren
      </button>
    </form>
    {#if error}
      <div class="mt-4 text-sm text-red-600 bg-red-100 rounded px-3 py-2">{error}</div>
    {/if}
  </div>
</div>
