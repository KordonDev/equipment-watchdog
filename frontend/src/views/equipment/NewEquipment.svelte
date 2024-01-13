<script lang="ts">
  import {onMount} from 'svelte';
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { routes } from "../../routes";
  import { push } from "svelte-spa-router";
  import {
    createEquipment,
    EquipmentType,
    translatedEquipmentTypes,
    type Equipment,
  } from "./equipment.service";
  import { getRegistrationCode } from "../registrationCode/registrationCode.service"
  import Navigation from "../../components/Navigation/Navigation.svelte";

  export let params = { type: undefined };

  let equipment: Equipment = {
    registrationCode: "",
    type: EquipmentType.Jacket,
    id: 0,
    size: "",
  };
  let loading = false;

  onMount(() => {
    if (params.type !== undefined) {
      equipment.type = params.type;
    }
    getRegistrationCode()
      .then(rc => {
        if (equipment.registrationCode === "") {
          equipment.registrationCode = rc.id
        }
      })
  })

  function handleSubmit() {
    loading = true;
    createEquipment(equipment)
      .then((e) => {
        createNotification(
          {
            color: "green",
            text: `Kleidungsstück angelegt.`,
          },
          5
        );
        loading = false;
        push(`${routes.EquipmentDetails.link}${e.id}`);
      })
      .catch(() => {
        createNotification({
          color: "red",
          text: `Kleidungsstück konnte nicht angelegt werden.`,
        });
        loading = false;
      });
  }
</script>

<Navigation />

<h1>Neues Kleidungsstück</h1>
<form on:submit|preventDefault={handleSubmit} disabled={loading}>
 <Label for="name" class="block mb-2">Registrierungsnummer</Label>
  <Input
    required
    class="mb-4"
    id="name"
    bind:value={equipment.registrationCode}
  />

  <Label class="mb-4">
    <div class="mb-2">Art</div>
    <Select
      required
      items={translatedEquipmentTypes}
      bind:value={equipment.type}
    />
  </Label>

  <Label for="size" class="block mb-2">Größe</Label>
  <Input required class="mb-4" id="size" bind:value={equipment.size} />
  <div class="flex flex-row justify-end">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="purple" type="submit">Anlegen</Button>
    {/if}
  </div>
</form>
