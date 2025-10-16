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

export function createEquipment(equipment: Equipment): Promise<Equipment> {
  return fetchApi("/equipment/", {
    method: "POST",
    body: JSON.stringify({
      ...equipment,
      id: 0,
    }),
  })
}

export function deleteEquipment(id: number): Promise<void> {
  return fetchApi(`/equipment/${id}`, {
    method: "DELETE",
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
