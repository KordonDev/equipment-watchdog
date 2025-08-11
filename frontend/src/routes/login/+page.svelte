<script lang="ts">
import type { PageData } from './$types';

interface Props {
  data: PageData;
}

let { data }: Props = $props();

let username = $state('');
let storeUsername = $state(false);
let loading = $state(false);

const handleLogin = (event: SubmitEvent) => {
  // Handle login logic
};

console.log(data.sections);
</script>

<div class="min-h-screen flex items-center justify-center bg-gray-100">
  <div class="w-full max-w-md bg-white rounded-lg shadow-lg p-8">
    <h2 class="text-2xl font-bold text-center text-gray-800 mb-6">Login</h2>
    <form onsubmit={handleLogin} class="space-y-4">
      <div>
        <label for="username" class="block mb-2 text-sm font-medium text-gray-700">Benutzer:</label>
        <input
          required
          class="block w-full rounded border border-gray-300 px-3 py-2 focus:border-purple-500 focus:ring focus:ring-purple-200 focus:ring-opacity-50"
          id="username"
          bind:value={username}
          type="text"
        />
      </div>
      <div class="flex items-center">
        <input
          type="checkbox"
          bind:checked={storeUsername}
          id="storeUsername"
          class="mr-2 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
        />
        <label for="storeUsername" class="text-sm text-gray-700 select-none">Benutzername speichern</label>
      </div>
      <button
        type="submit"
        class="w-full rounded bg-purple-600 px-4 py-2 text-white font-semibold hover:bg-purple-700 transition-colors duration-150 disabled:opacity-50"
        disabled={loading}
      >
        Einloggen
      </button>
    </form>
  </div>
</div>
