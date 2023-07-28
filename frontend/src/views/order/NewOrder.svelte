<script lang="ts">
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import { createNotification } from "../../components/Notification/notificationStore";
  import { routes } from "../../routes";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte/internal";

  import {
    EquipmentType,
    translatedEquipmentTypes,
    type Equipment,
  } from "../equipment/equipment.service";
  import {
    createOrder,
  } from "./order.service";
  import {
    getMember
  } from '../member/member.service'
  import Navigation from "../../components/Navigation/Navigation.svelte";

  export let params = { memberId: undefined };
  const memberId = parseInt(params.memberId);

  let order: Order = {
    type: EquipmentType.Jacket,
    id: 0,
    size: "",
    memberId: memberId,
  };
  let loading = false;
  let memberName = "";
  
  onMount(() => {
    getMember(memberId)
      .then((member) => {
        memberName = member.name;
        console.log(memberName)
      });
  })



  function handleSubmit() {
    loading = true;
    createOrder(order)
      .then((e) => {
        createNotification(
          {
            color: "green",
            text: `Bestellung angelegt.`,
          },
          5
        );
        loading = false;
        push(`${routes.OrderDetails.link}${e.id}`);
      })
      .catch(() => {
        createNotification({
          color: "red",
          text: `Bestellung konnte nicht angelegt werden.`,
        });
        loading = false;
      });
  }
</script>

<Navigation />

<h1>Neue Bestellung</h1>
<form on:submit|preventDefault={handleSubmit} disabled={loading}>
  <Label for="name" class="block mb-2">Für</Label>
  <Input
    required
    class="mb-4"
    id="name"
    disabled
    bind:value={memberName}
  />

  <Label class="mb-4">
    <div class="mb-2">Art</div>
    <Select
      required
      items={translatedEquipmentTypes}
      bind:value={order.type}
    />
  </Label>

  <Label for="size" class="block mb-2">Größe</Label>
  <Input required class="mb-4" id="size" bind:value={order.size} />
  <div class="flex flex-row justify-end">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="purple" type="submit">Anlegen</Button>
    {/if}
  </div>
</form>
