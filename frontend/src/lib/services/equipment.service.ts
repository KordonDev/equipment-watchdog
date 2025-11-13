import { fetchApi } from "../apiService";
import { Group } from '$lib/services/member.service';

export interface Equipment {
  id: number;
  type: EquipmentType;
  registrationCode: string;
  memberId?: number;
  size: string;
}


export function getEquipment(id: string): Promise<Equipment> {
  return fetchApi(`/equipment/${id}`);
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

export function randomRegistrationCode(): string {
	return (Math.random() + 1).toString(36).substring(7)
}

export type EquipmentListByType = {
  [EquipmentType.Helmet]?: Equipment[];
  [EquipmentType.Jacket]?: Equipment[];
  [EquipmentType.Gloves]?: Equipment[];
  [EquipmentType.Trousers]?: Equipment[];
  [EquipmentType.Boots]?: Equipment[];
  [EquipmentType.TShirt]?: Equipment[];
};

export function getFreeEquipment(): Promise<EquipmentListByType> {
  return fetchApi("/equipment/free");
}

export function getAllEquipment(): Promise<Equipment[]> {
	return fetchApi("/equipment/");
}

export enum EquipmentType {
  Helmet = "helmet",
  Jacket = "jacket",
  Gloves = "gloves",
  Trousers = "trousers",
  Boots = "boots",
  TShirt = "tshirt",
}

export const equipmentLabels = {
	[EquipmentType.Helmet]: 'Helm',
	[EquipmentType.Jacket]: 'Jacke',
	[EquipmentType.Gloves]: 'Handschuhe',
	[EquipmentType.Trousers]: 'Hose',
	[EquipmentType.Boots]: 'Stiefel',
	[EquipmentType.TShirt]: 'T-Shirt'
};

export const translateEquipmentType = (type: EquipmentType) => equipmentLabels[type] || type;
