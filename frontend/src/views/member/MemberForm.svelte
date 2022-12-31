<script lang="ts">
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { getTranslatedGroups } from "../../components/groupsStore";
  import {
    addNoneEquipment,
    allEquipmentPossibilities,
    EquipmentType,
    getEquipmentFromFree,
    getFreeEquipment,
  } from "../equipment/equipment.service";
  import { Group, type Member } from "./member.service";

  export let member: Member;
  export let submitText: string;
  export let loading: boolean;
  export let onSubmit: (m: Member) => void;
  export let hideEquipment: boolean | undefined;

  let allGroups = [];
  getTranslatedGroups.subscribe((groups) => (allGroups = groups));
  let mocked = false;

  onMount(() => {
    member.equipments = addNoneEquipment(member.equipments);
    mocked = true;
  });

  let freeEquipment;
  getFreeEquipment().then((fe) => {
    freeEquipment = allEquipmentPossibilities(member.equipments, fe);
  });

  function handleSubmit() {
    onSubmit({
      ...member,
      equipments: getEquipmentFromFree(member.equipments, freeEquipment),
    });
  }

  function getItems(
    freeEquipment: EquipmentType,
    equipmentType: EquipmentType
  ) {
    return (
      freeEquipment[equipmentType]?.map((e) => ({
        value: e.registrationCode,
        name: e.registrationCode,
      })) || []
    );
  }
</script>

<form on:submit|preventDefault={handleSubmit} disabled={loading}>
  <Label for="name" class="block mb-2">Name</Label>
  <Input required class="mb-4" id="name" bind:value={member.name} />

  <Label class="mb-4">
    <div class="mb-2">Gruppe</div>
    <Select required items={allGroups} bind:value={member.group} />
  </Label>
  {#if mocked && freeEquipment && !hideEquipment}
    <div>
      <h2>Ausrüstung</h2>
      <Label class="mb-4">
        <div class="mb-2">Helm</div>
        <Select
          items={getItems(freeEquipment, EquipmentType.Helmet)}
          bind:value={member.equipments[EquipmentType.Helmet].registrationCode}
        />
      </Label>
      {#if member.group !== Group.MINI}
        <Label class="mb-4">
          <div class="mb-2">Jacke</div>
          <Select
            items={getItems(freeEquipment, EquipmentType.Jacket)}
            bind:value={member.equipments[EquipmentType.Jacket]
              .registrationCode}
          />
        </Label>
      {/if}
      <Label class="mb-4">
        <div class="mb-2">Handschuhe</div>
        <Select
          items={getItems(freeEquipment, EquipmentType.Gloves)}
          bind:value={member.equipments[EquipmentType.Gloves].registrationCode}
        />
      </Label>
      {#if member.group !== Group.MINI}
        <Label class="mb-4">
          <div class="mb-2">Hose</div>
          <Select
            items={getItems(freeEquipment, EquipmentType.Trousers)}
            bind:value={member.equipments[EquipmentType.Trousers]
              .registrationCode}
          />
        </Label>
        <Label class="mb-4">
          <div class="mb-2">Stiefel</div>
          <Select
            items={getItems(freeEquipment, EquipmentType.Boots)}
            bind:value={member.equipments[EquipmentType.Boots].registrationCode}
          />
        </Label>
        <Label class="mb-4">
          <div class="mb-2">TShirt</div>
          <Select
            items={getItems(freeEquipment, EquipmentType.TShirt)}
            bind:value={member.equipments[EquipmentType.TShirt]
              .registrationCode}
          />
        </Label>
      {/if}
    </div>
  {/if}
  <div class="flex flex-row justify-end">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="purple" type="submit">
        {submitText}
      </Button>
    {/if}
  </div>
</form>
