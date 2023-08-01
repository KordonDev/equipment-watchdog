<script lang="ts">
  import Navigation from "../../components/Navigation/Navigation.svelte";
  import { deleteOrder, getOrder, fulfillOrder } from "./order.service";
  import { Alert,Label, Input, Button, Modal, Spinner } from "flowbite-svelte";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { push } from "svelte-spa-router";
  import { routes } from "../../routes";
  import { translateEquipmentType } from "../equipment/equipment.service";

  export let order = {};

  let fulfillModalOpen = false;
  let loading = false;
  let registrationCode = "";

  console.log(order)

  function fulfillOrder_internal() {
    loading = true;
    fulfillOrder(order, registrationCode)
      .then((equipment) => {
        createNotification(
          {
            color: "green",
            text: `Bestellung wurde zu einem Ausrüstungsgegenstand.`,
          },
          5
        );
        loading = false;
        fulfillModalOpen = false;
        push(`${routes.MemberDetail.link}${equipment.memberId}`);
      })
      .catch(() => {
        createNotification(
          {
            color: "red",
            text: `Bestellung konnte nicht zu einem Ausrüstungsgegenstand gemacht und zugewiesen werden. Prüfe, ob es einen neuen Ausrüstungsgegenstand gibt und probiere sonst nocheinmal.`,
          },
          20
        );
        fulfillModalOpen = false;
        loading = false;
      });
  }
</script>

<Button color="green" on:click={() => (fulfillModalOpen = true)}>
  Bestellung einlösen
</Button>
<Modal
  title={`Bestellung ${translateEquipmentType(order.type)} einlösen`}
  bind:open={fulfillModalOpen}
>
  <p>
    Soll die Bestellung für {translateEquipmentType(order.type)} eingelöst werden?
  </p>
    <form on:submit|preventDefault={() => fulfillOrder_internal()} disabled={loading} id="fulfilOrderForm">
      <Label for="registrationCode" class="block mb-2">Registrierungsnummer</Label>
      <Input
        required
        class="mb-4"
        id="registrationCode"
        bind:value={registrationCode}
      />
    </form>
  <svelte:fragment slot="footer">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="green" form="fulfilOrderForm" type="submit">
        Einlösen
      </Button>
      <Button color="purple">Abbrechen</Button>
    {/if}
  </svelte:fragment>
</Modal>
