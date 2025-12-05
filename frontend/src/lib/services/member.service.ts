import { fetchApi } from "../apiService";
import  type { Equipment } from "./equipment.service";
import {EquipmentType} from './equipment.service';

export interface Member {
  id: number;
  name: string;
  group: string;
  equipments: EquipmentByType;
}

export interface EquipmentByType {
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

export function saveEquipmentForMember(
	memberId: number,
	equipmentType: EquipmentType,
	equipment: Equipment
): Promise<Equipment> {
	return fetchApi(`/members/${memberId}/${equipmentType}`, {
		method: "POST",
		body: JSON.stringify({
			registrationCode: equipment.registrationCode,
			size: equipment.size,
		}),
	});
}

export function removeEquipmentFromMember(
	memberId: number,
	equipmentType: EquipmentType,
): Promise<Equipment> {
	return fetchApi(`/members/${memberId}/${equipmentType}`, {
		method: "DELETE",
	});
}


export enum Group {
	MONDAY = "monday",
  FRIDAY = "friday",
  MINI = "mini",
}

export const groupLabels = {
	[Group.FRIDAY]: 'Freitag',
	[Group.MINI]: 'Mini',
	[Group.MONDAY]: 'Montag'
};

export const equipmentForGroup: any = {
  [Group.FRIDAY]: [EquipmentType.Helmet, EquipmentType.Jacket, EquipmentType.Gloves, EquipmentType.Trousers, EquipmentType.Boots, EquipmentType.TShirt],
	[Group.MONDAY]: [EquipmentType.Helmet, EquipmentType.Jacket, EquipmentType.Gloves, EquipmentType.Trousers, EquipmentType.Boots, EquipmentType.TShirt],
  [Group.MINI]: [EquipmentType.Helmet, EquipmentType.Gloves, EquipmentType.TShirt]
}
