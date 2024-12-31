import { map, onMount } from "nanostores";
import { GetData, UpdateData } from "../../wailsjs/go/main/App";
import type { main } from "../../wailsjs/go/models";

export const $app = map<main.AppData>({
  macros: [],
});

onMount($app, () => {
  GetData().then(data => $app.set(data));
});

export async function saveMacros(macros: main.Macro[]) {
  const data: main.AppData = { macros };
  await UpdateData(data);
  $app.set(data);
}

export async function addMacro(macro: main.Macro) {
  const { macros } = $app.get();
  await saveMacros([...macros, macro]);
}

export async function updateMacro(updatedMacro: main.Macro) {
  const { macros } = $app.get();
  const updatedMacros = macros.map(macro => (macro.id === updatedMacro.id ? updatedMacro : macro));
  await saveMacros(updatedMacros);
}

export async function deleteMacro(macroId: string) {
  const { macros } = $app.get();
  const filteredMacros = macros.filter(macro => macro.id !== macroId);
  await saveMacros(filteredMacros);
}
