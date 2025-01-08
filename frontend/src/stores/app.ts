import { map, onMount } from "nanostores";
import {
  ReadAppData,
  WriteAppData,
} from "../../bindings/repeat-what-shit/internal/app";
import { AppData, Macro } from "../../bindings/repeat-what-shit/internal/types";

export const $app = map<AppData>({
  macros: [],
});

onMount($app, () => {
  ReadAppData().then((data) => $app.set(data));
});

export async function saveMacros(macros: Macro[]) {
  const data: AppData = { macros };
  await WriteAppData(data);
  $app.set(data);
}

export async function addMacro(macro: Macro) {
  const { macros } = $app.get();
  await saveMacros([...(macros || []), macro]);
}

export async function updateMacro(updatedMacro: Macro) {
  const { macros } = $app.get();
  const updatedMacros = macros?.map((macro) =>
    macro.id === updatedMacro.id ? updatedMacro : macro
  );
  await saveMacros(updatedMacros || []);
}

export async function deleteMacro(macroId: string) {
  const { macros } = $app.get();
  const filteredMacros = macros?.filter((macro) => macro.id !== macroId);
  await saveMacros(filteredMacros || []);
}

export async function toggleMacroDisabled(macroId: string) {
  const { macros } = $app.get();
  const updatedMacros = macros?.map((macro) =>
    macro.id === macroId ? { ...macro, disabled: !macro.disabled } : macro
  );
  await saveMacros(updatedMacros || []);
}
