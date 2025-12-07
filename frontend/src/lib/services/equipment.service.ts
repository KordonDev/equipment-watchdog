import { fetchApi } from "../apiService";

export interface Equipment {
  id: number;
  type: EquipmentType;
  registrationCode: string;
  memberId?: number;
  size: string;
}

export function randomRegistrationCode(): string {
	return (Math.random() + 1).toString(36).substring(7)
}

export type EquipmentListByType = {
  [EquipmentType.Helmet]?: Equipment[];
  [EquipmentType.Jacket]?: Equipment[];
	[EquipmentType.OverJacket]?: Equipment[];
  [EquipmentType.Gloves]?: Equipment[];
  [EquipmentType.Trousers]?: Equipment[];
  [EquipmentType.Boots]?: Equipment[];
  [EquipmentType.TShirt]?: Equipment[];
};

export function getAllEquipment(): Promise<Equipment[]> {
	return fetchApi("/equipment/");
}

export enum EquipmentType {
  Helmet = "helmet",
  Jacket = "jacket",
	OverJacket = "overJacket",
  Gloves = "gloves",
  Trousers = "trousers",
  Boots = "boots",
  TShirt = "tshirt",
}

export const equipmentLabels = {
	[EquipmentType.Helmet]: 'Helm',
	[EquipmentType.Jacket]: 'Jacke',
	[EquipmentType.OverJacket]: 'Überjacke',
	[EquipmentType.Gloves]: 'Handschuhe',
	[EquipmentType.Trousers]: 'Hose',
	[EquipmentType.Boots]: 'Stiefel',
	[EquipmentType.TShirt]: 'T-Shirt'
};

export const translateEquipmentType = (type: EquipmentType) => equipmentLabels[type] || type;
