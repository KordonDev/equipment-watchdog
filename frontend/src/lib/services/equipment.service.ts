import { fetchApi } from "../apiService";
import type { EquipmentByType } from "./member.service";
import { useRegistrationCode } from "./registrationCode.service";

export interface Equipment {
  id: number;
  type: EquipmentType;
  registrationCode: string;
  memberId?: number;
  size: string;
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
  })
    .then(e => {
      useRegistrationCode(e.registrationCode);
      return e;
    });
}

export function deleteEquipment(id: number): Promise<void> {
  return fetchApi(`/equipment/${id}`, {
    method: "DELETE",
  });
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

export const NoEquipment: Equipment = {
  id: 0,
  registrationCode: "Keine",
  type: EquipmentType.TShirt,
  size: ""
};

export function addNoneEquipment(equipments: EquipmentByType): EquipmentByType {
  return (Object.values(EquipmentType) as EquipmentType[]).reduce(
    (acc: EquipmentByType, eT: EquipmentType) => {
      if (equipments[eT]) {
        acc[eT] = equipments[eT];
      } else {
        acc[eT] = { ...NoEquipment };
      }
      return acc;
    },
    {} as EquipmentByType
  );
}

export function getEquipmentFromFree(
  selected: EquipmentByType,
  equipments: EquipmentListByType
): EquipmentByType {
  return (Object.values(EquipmentType) as EquipmentType[]).reduce((acc, eT) => {
    const selectedEquipment = equipments[eT]?.find(
      (e) => selected[eT]?.registrationCode === e.registrationCode
    );
    if (
      selected[eT]?.registrationCode !== NoEquipment.registrationCode &&
      selectedEquipment
    ) {
      acc[eT] = selectedEquipment;
    }
    return acc;
  }, {} as EquipmentByType);
}

export function allEquipmentPossibilities(
  equipments: EquipmentByType,
  equipmentsList: EquipmentListByType
): EquipmentListByType {
  return (Object.values(EquipmentType) as EquipmentType[]).reduce((acc, eT) => {
    acc[eT] = [NoEquipment];
    if (
      equipments[eT] &&
      equipments[eT].registrationCode !== NoEquipment.registrationCode
    ) {
      acc[eT] = acc[eT].concat(equipments[eT]);
    }
    if (equipmentsList[eT]) {
      acc[eT] = acc[eT].concat(...equipmentsList[eT]);
    }
    return acc;
  }, {} as EquipmentListByType);
}
