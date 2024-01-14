<script lang="ts">
  import { Label, Input, Select, Button, Spinner } from "flowbite-svelte";
  import EquipmentSelect from "../../components/Equipment/EquipmentSelect.svelte";
  import type { Writable } from "svelte/store";
  import { getTranslatedGroups } from "../../components/groupsStore";
  import {
    allEquipmentPossibilities,
    EquipmentType,
    getFreeEquipment,
  } from "../equipment/equipment.service";
  import { Group, type Member } from "./member.service";

  export let memberStore: Writable<Member>;
  export let submitText: string;
  export let loading: boolean;
  export let onSubmit: (m: Member) => void;
  export let hideEquipment: boolean | undefined;

  let allGroups = [];
  getTranslatedGroups.subscribe((groups) => (allGroups = groups));

  let freeEquipment = allEquipmentPossibilities($memberStore.equipments, {});
  getFreeEquipment().then((fe) => {
    freeEquipment = allEquipmentPossibilities($memberStore.equipments, fe);
  });

  function handleSubmit() {
    onSubmit($memberStore);
  }
</script>

<form on:submit|preventDefault={handleSubmit}>
  <Label for="name" class="block mb-2">Name</Label>
  <Input required class="mb-4" id="name" bind:value={$memberStore.name} />

  <Label class="mb-4">
    <div class="mb-2">Gruppe</div>
    <Select required items={allGroups} bind:value={$memberStore.group} />
  </Label>
  {#if freeEquipment && !hideEquipment}
    <div>
      <h2>Ausrüstung</h2>
      <EquipmentSelect
        freeEquipments={freeEquipment}
        equipmmentType={EquipmentType.Helmet}
        {memberStore}
      />
      {#if $memberStore.group !== Group.MINI}
        <EquipmentSelect
          freeEquipments={freeEquipment}
          equipmmentType={EquipmentType.Jacket}
          {memberStore}
        />
      {/if}
      <EquipmentSelect
        freeEquipments={freeEquipment}
        equipmmentType={EquipmentType.Gloves}
        {memberStore}
      />
      {#if $memberStore.group !== Group.MINI}
        <EquipmentSelect
          freeEquipments={freeEquipment}
          equipmmentType={EquipmentType.Trousers}
          {memberStore}
        />
        <EquipmentSelect
          freeEquipments={freeEquipment}
          equipmmentType={EquipmentType.Boots}
          {memberStore}
        />
      {/if}
      <EquipmentSelect
        freeEquipments={freeEquipment}
        equipmmentType={EquipmentType.TShirt}
        {memberStore}
      />
    </div>
  {/if}
  <div class="flex flex-row justify-end">
    {#if loading}
      <Spinner />
    {:else}
      <Button color="purple" type="submit" disabled={loading}>
        {submitText}
      </Button>
    {/if}
  </div>
</form>
