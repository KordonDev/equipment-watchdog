<script lang="ts">
  import { Group, type Member } from "../../views/member/member.service";
  import { EquipmentType } from "../../views/equipment/equipment.service";
  import EquipmentIcon from "./EquipmentIcon.svelte";

  export let member: Member;
  export let columns: number;
</script>

<div class={`flex flex-wrap columns-${columns} align-center`}>
  <EquipmentIcon
    registrationCode={member.equipments.helmet?.registrationCode}
    equipmentType={EquipmentType.Helmet}
  />
  {#if member.group !== Group.MINI}
    <EquipmentIcon
      registrationCode={member.equipments.jacket?.registrationCode}
      equipmentType={EquipmentType.Jacket}
    />
  {/if}
  <EquipmentIcon
    registrationCode={member.equipments.gloves?.registrationCode}
    equipmentType={EquipmentType.Gloves}
  />
  {#if member.group !== Group.MINI}
    <EquipmentIcon
      registrationCode={member.equipments.trousers?.registrationCode}
      equipmentType={EquipmentType.Trousers}
    />
    <EquipmentIcon
      registrationCode={member.equipments.boots?.registrationCode}
      equipmentType={EquipmentType.Boots}
    />
  {/if}
  <EquipmentIcon
    registrationCode={member.equipments.tshirt?.registrationCode}
    equipmentType={EquipmentType.TShirt}
  />
</div>

<style>
  .columns-2 {
    max-width: 260px;
  }
  .columns-6 {
    max-width: 780px;
  }
</style>
