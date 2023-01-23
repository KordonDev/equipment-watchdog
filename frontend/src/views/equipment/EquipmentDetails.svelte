<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import {
    deleteEquipment,
    EquipmentType,
    getEquipment,
    translateEquipmentType,
  } from "./equipment.service";
  import { Alert, Button, Modal, Spinner } from "flowbite-svelte";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { push } from "svelte-spa-router";
  import { routes } from "../../routes";

  export let params = { id: undefined };

  let deleteModalOpen = false;
  let loading = false;
  let equipmentPromise = getEquipment(params.id);

  function deleteEquipmentInternal(id: number, type: EquipmentType) {
    deleteEquipment(id)
      .then(() => {
        createNotification(
          {
            color: "green",
            text: `${translateEquipmentType(type)} wurde erfolgreich gelöscht.`,
          },
          5
        );
        loading = false;
        deleteModalOpen = false;
        push(`${routes.EquipmentType.link}${type}`);
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `${translateEquipmentType(
              type
            )} konnte nicht gelöscht werden.`,
          },
          20
        );
        loading = false;
      });
  }
</script>

<Navigation />

{#await equipmentPromise}
  <Spinner />
{:then equipment}
  <div>
    <h3>{translateEquipmentType(equipment.type)}</h3>
    <p>Registrierungsnummer {equipment.registrationCode}</p>
    <p>Größe {equipment.size || "-"}</p>
  </div>
  <Button color="red" on:click={() => (deleteModalOpen = true)}>
    {translateEquipmentType(equipment.type)} löschen
  </Button>
  <Modal
    title={`${translateEquipmentType(equipment.type)} löschen`}
    bind:open={deleteModalOpen}
    autoclose
  >
    <p>
      Soll {translateEquipmentType(equipment.type)} mit der Registrierungsnummer
      {equipment.registrationCode} gelöscht werden?
    </p>
    <svelte:fragment slot="footer">
      {#if loading}
        <Spinner />
      {:else}
        <Button
          color="red"
          on:click={() => deleteEquipmentInternal(equipment.id, equipment.type)}
        >
          Löschen
        </Button>
        <Button color="purple">Abbrechen</Button>
      {/if}
    </svelte:fragment>
  </Modal>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnte die Ausrüstung nicht geladen werden.
  </Alert>
{/await}
