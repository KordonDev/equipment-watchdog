import { readable } from "svelte/store";
import type { Groups } from "../views/member/member.service";
import {
  getGroups as getGroupsRequest,
  Group,
} from "../views/member/member.service";

export const getGroupsWithEquipment = readable<Groups>(
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

export const getTranslatedGroups = readable<TranslatedGroup[]>(
  undefined,
  function start(set) {
    getGroupsRequest().then((groups) => {
      const translations = Object.keys(groups).map((key) => {
        const translation = groupTranslations.find(
          (t) => t.value === key
        )?.translation;
        return {
          value: key as Group,
          name: translation || key,
        };
      });

      set(translations);
    });
  }
);
