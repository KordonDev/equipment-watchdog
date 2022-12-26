<script lang="ts">
  import type { Member } from "../../views/member/member.service";
  import { EquipmentType } from "../../views/equipment/equipment.service";
  import EquipmentIcon from "./EquipmentIcon.svelte";

  export let member: Member;
  export let columns: number;
</script>

<div class={`flex flex-wrap columns-${columns}`}>
  <EquipmentIcon
    registrationCode={member.equipments.helmet?.registrationCode}
    equipmentType={EquipmentType.Helmet}
  />
  <EquipmentIcon
    registrationCode={member.equipments.jacket?.registrationCode}
    equipmentType={EquipmentType.Jacket}
  />
  <EquipmentIcon
    registrationCode={member.equipments.gloves?.registrationCode}
    equipmentType={EquipmentType.Gloves}
  />
  <EquipmentIcon
    registrationCode={member.equipments.trousers?.registrationCode}
    equipmentType={EquipmentType.Trousers}
  />
  <EquipmentIcon
    registrationCode={member.equipments.boots?.registrationCode}
    equipmentType={EquipmentType.Boots}
  />
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
