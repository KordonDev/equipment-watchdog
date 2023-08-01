<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import FulfillOrderModal from './FulfillOrderModal.svelte';
  import { deleteOrder, getOrder, fulfillOrder } from "./order.service";
  import { Alert, Button, Modal, Spinner } from "flowbite-svelte";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { push } from "svelte-spa-router";
  import { routes } from "../../routes";
  import { translateEquipmentType } from "../equipment/equipment.service";

  export let params = { id: undefined };

  let deleteModalOpen = false;
  let loading = false;
  let orderPromise = getOrder(params.id);

  function deleteOrderInternal(id: number) {
    deleteOrder(id)
      .then(() => {
        createNotification(
          {
            color: "green",
            text: `Bestellung ${id} wurde erfolgreich gelöscht.`,
          },
          5
        );
        loading = false;
        deleteModalOpen = false;
        push(routes.Orders.link);
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Bestellung konnte nicht gelöscht werden.`,
          },
          20
        );
        loading = false;
      });
  }
</script>

<Navigation />

{#await orderPromise}
  <Spinner />
{:then order}
  <div>
    <h3>{translateEquipmentType(order.type)}</h3>
    <p>Bestell Id {order.id}</p>
    <p>Bestelldatum {order.createdAt}</p>
    <p>Lieferdatum {order.fulfilledAt || "-"}</p>
    <p>Größe {order.size || "-"}</p>
  </div>
  <FulfillOrderModal order={order}/>
  <Button color="red" on:click={() => (deleteModalOpen = true)}>
    Bestellung löschen
  </Button>
  <Modal
    title={`Bestellung ${translateEquipmentType(order.type)} löschen`}
    bind:open={deleteModalOpen}
    autoclose
  >
    <p>
      Soll die Bestellung für {translateEquipmentType(order.type)} gelöscht werden?
    </p>
    <svelte:fragment slot="footer">
      {#if loading}
        <Spinner />
      {:else}
        <Button color="red" on:click={() => deleteOrderInternal(order.id)}>
          Löschen
        </Button>
        <Button color="purple">Abbrechen</Button>
      {/if}
    </svelte:fragment>
  </Modal>
{:catch}
  <Alert class="mb-4 top-alert" color="red" dismissable style="display: block;">
    Leider konnte die Bestellung nicht geladen werden.
  </Alert>
{/await}
