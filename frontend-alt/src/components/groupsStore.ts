import { writable } from "svelte/store";
import type { Groups } from "../views/member/member.service";
import {
  getGroups as getGroupsRequest,
  Group,
} from "../views/member/member.service";

export const getGroupsWithEquipment = writable<Groups>(
  undefined,
  function start(set) {
    getGroupsRequest().then((groups) => set(groups));
  }
);

export interface TranslatedGroup {
  name: string;
  value: Group;
}

const groupTranslations = [
  { value: Group.MONDAY, translation: "Montagsgruppe" },
  { value: Group.FRIDAY, translation: "Freitagsgruppe" },
  { value: Group.MINI, translation: "Minigruppe" },
];

export const getTranslatedGroups = writable<TranslatedGroup[]>(
  undefined,
  function start(set) {
    loadTranslatedGroups().then((translations) => {
      set(translations);
    })
  }
);

export const loadTranslatedGroups = async () => {
  const groups = await getGroupsRequest();
  if (!groups) {
    return;
  }
  return Object.keys(groups).map((key) => {
    const translation = groupTranslations.find(
        (t) => t.value === key
    )?.translation;
    return {
      value: key as Group,
      name: translation || key,
    };
  });
}

const LOCALSTOARGE_GROUP_KEY = "groups_overview";
const storeGroup = localStorage.getItem(LOCALSTOARGE_GROUP_KEY);
export const groupFilter = writable(storeGroup || Group.MONDAY);
groupFilter.subscribe((value) =>
  localStorage.setItem(LOCALSTOARGE_GROUP_KEY, value)
);
