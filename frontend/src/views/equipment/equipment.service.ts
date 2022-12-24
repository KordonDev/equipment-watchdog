import { fetchApi } from "../apiService";

export interface Equipment {
  id: number;
  type: EquipmentType;
  registrationCode: string;
}

export function getEquipmentByType(type: EquipmentType): Promise<Equipment[]> {
  return fetchApi(`/equipment/type/${type}`);
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
  });
}

export function deleteEquipment(id: number): Promise<void> {
  return fetchApi(`/equipment/${id}`, {
    method: "DELETE",
  });
}

export type EquipmentByType = {
  [EquipmentType.Helmet]?: Equipment[];
  [EquipmentType.Jacket]?: Equipment[];
  [EquipmentType.Gloves]?: Equipment[];
  [EquipmentType.Trousers]?: Equipment[];
  [EquipmentType.Boots]?: Equipment[];
  [EquipmentType.TShirt]?: Equipment[];
};

export function getFreeEquipment(): Promise<EquipmentByType> {
  return fetchApi("/equipment/free");
}

export enum EquipmentType {
  Helmet = "helmet",
  Jacket = "jacket",
  Gloves = "gloves",
  Trousers = "trousers",
  Boots = "boots",
  TShirt = "tshirt",
}

export const translatedEquipmentTypes = [
  { value: EquipmentType.Helmet, name: "Helm" },
  { value: EquipmentType.Jacket, name: "Jacke" },
  { value: EquipmentType.Gloves, name: "Handschue" },
  { value: EquipmentType.Trousers, name: "Hose" },
  { value: EquipmentType.Boots, name: "Stiefel" },
  { value: EquipmentType.TShirt, name: "T-Shirt" },
];

export const translateEquipmentType = (type: EquipmentType): string => {
  return (
    translatedEquipmentTypes.find(
      (equipmentType) => equipmentType.value === type
    )?.name || type
  );
};
