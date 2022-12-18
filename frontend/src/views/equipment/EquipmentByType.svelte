<script lang="ts">
  import { routes } from "../../routes";
  import { link, push } from "svelte-spa-router";
  import { afterUpdate, onMount } from "svelte/internal";
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import {
    EquipmentType,
    getEquipmentByType,
    translatedEquipmentTypes,
    translateEquipmentType,
  } from "./equipment.service";
  import { Alert, Label, Select, Spinner } from "flowbite-svelte";

  export let params = { type: undefined };
  let currentType;

  onMount(() => {
    currentType = params.type;
    if (translateEquipmentType(params.type) === params.type) {
      push(`${routes.EquipmentType.link}${EquipmentType.Jacket}`);
    }
  });

  afterUpdate(() => {
    if (currentType !== params.type) {
      currentType = params.type;
      equipmentsPromise = getEquipmentByType(params.type);
    }
  });

  function updateUrl() {
    push(`${routes.EquipmentType.link}${params.type}`);
  }

  let equipmentsPromise = getEquipmentByType(params.type);
</script>

<Navigation />
<h1>Ausrüstung vom Typ {translateEquipmentType(params.type)}</h1>

<a href={routes.EquipmentAdd.link} use:link>Neue Ausrüstung anlegen</a>

<Label class="mb-4">
  <div class="mb-2">Art</div>
  <Select
    required
    items={translatedEquipmentTypes}
    bind:value={params.type}
    on:change={updateUrl}
  />
</Label>
{#await equipmentsPromise}
  <Spinner />
{:then equipments}
  <div class="flex flex-wrap">
    {#each equipments as equipment}
      <div>
        <h3>{translateEquipmentType(equipment.type)}</h3>
        <p>{equipment.type}</p>
      </div>
    {/each}
  </div>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnte die Ausrüstung nicht geladen werden.
  </Alert>
{/await}
