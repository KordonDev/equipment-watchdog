import { fetchApi } from "../apiService";
import type { Equipment, EquipmentType } from "../equipment/equipment.service";
import { EquipmentType as EquipmentTypes } from "../equipment/equipment.service";

export interface Member {
  id: string;
  name: string;
  group: string;
  equipments: Equipments;
}

export interface Equipments {
  [EquipmentType.Helmet]?: Equipment;
  [EquipmentType.Jacket]?: Equipment;
  [EquipmentType.Gloves]?: Equipment;
  [EquipmentType.Trousers]?: Equipment;
  [EquipmentType.Boots]?: Equipment;
  [EquipmentType.TShirt]?: Equipment;
}

export function getMembers(): Promise<Member[]> {
  return fetchApi("/members/");
}

export function getMember(id: string): Promise<Member> {
  return fetchApi(`/members/${id}`);
}

export function createMember(member: Member): Promise<Member> {
  return fetchApi("/members/", {
    method: "POST",
    body: JSON.stringify({
      ...member,
      id: 0,
    }),
  });
}

export function updateMember(member: Member): Promise<Member> {
  return fetchApi(`/members/${member.id}`, {
    method: "PUT",
    body: JSON.stringify(member),
  });
}

export function deleteMember(id: string): Promise<void> {
  return fetchApi(`/members/${id}`, {
    method: "DELETE",
  });
}

export enum Group {
  FRIDAY = "friday",
  MONDAY = "monday",
  MINI = "mini",
}

export interface Groups {
  [Group.FRIDAY]: EquipmentType[];
  [Group.MONDAY]: EquipmentType[];
  [Group.MINI]: EquipmentType[];
}

export function getGroups(): Promise<Groups> {
  return fetchApi("/members/groups");
}
