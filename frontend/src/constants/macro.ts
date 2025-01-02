import { main } from "../../wailsjs/go/models";

export enum MacroType {
  Sequence = 0,
  Toggle = 1,
}

export const MACRO_TYPE_NAMES: Record<MacroType, string> = {
  [MacroType.Sequence]: "Обычный",
  [MacroType.Toggle]: "Переключатель",
};

export const DEFAULT_MACRO: main.Macro = {
  id: "",
  name: "",
  activation_keys: [],
  type: 0,
  actions: [],
  include_titles: "",
};
