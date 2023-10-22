<script lang="ts">
  import { Heading } from "flowbite-svelte";
  import { Alert, Spinner } from "flowbite-svelte";

  export let changesPromise: Promise<string>;
</script>

{#await changesPromise}
  <Spinner />
{:then changes}
  <div>
    <Heading tag="h3" class="mb-4">Letzte Änderungen</Heading>
    <ul>
      {#each changes as change}
        <li>{change}</li>
      {/each}
    </ul>
    {#if changes.length === 0}
      <p>Keine Änderungen.</p>
    {/if}
  </div>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Änderungen konnten nicht aufgerufen werden.
  </Alert>
{/await}
