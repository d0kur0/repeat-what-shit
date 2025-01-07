import { useNavigate } from "@solidjs/router";
import { ViewPort } from "../components/ViewPort";
import btnStyles from "../styles/Buttons.module.css";
import styles from "./MacrosList.module.css";
import { useStore } from "@nanostores/solid";
import { $app, toggleMacroDisabled } from "../stores/app";
import {
  Macro,
  MacroType,
} from "../../bindings/repeat-what-shit/internal/types";
import { For } from "solid-js";
import { getKeyName } from "../utils/keys";
import { Switch } from "@kobalte/core/switch";

import switchStyles from "../styles/Switch.module.css";

const macroTypeMap = {
  [MacroType.MacroTypeSequence]: "По нажатию",
  [MacroType.MacroTypeToggle]: "Переключение",
  [MacroType.MacroTypeHold]: "Удержание",
};

function Macros(props: { macros: Macro; onToggleDisable?: () => void }) {
  const navigate = useNavigate();

  return (
    <button
      onClick={() => navigate(`/macros/edit/${props.macros.id}`)}
      class={styles.macro}
    >
      <div class="flex gap-2 items-center">
        <Switch
          onChange={props.onToggleDisable}
          checked={!props.macros.disabled}
          onClick={(e: Event) => e.stopPropagation()}
          class={switchStyles.root}
        >
          <Switch.Input class={switchStyles.input} />
          <Switch.Control class={switchStyles.control}>
            <Switch.Thumb class={switchStyles.thumb} />
          </Switch.Control>
        </Switch>
        <div class={styles.macroName}>{props.macros.name}</div>
        <div class="flex-1" />
        <div class={styles.macroType}>{macroTypeMap[props.macros.type]}</div>
      </div>

      <div class={styles.keysTitle}>Клавиши активации</div>

      <div class="flex gap-2 mt-2">
        <For each={props.macros.activation_keys}>
          {(k) => <div class={styles.kbd}>{getKeyName(k)}</div>}
        </For>
      </div>

      {!!props.macros.include_title?.length || (
        <div class={styles.macroWindows}>Без ограничения по окнам</div>
      )}

      {!!props.macros.include_title?.length && (
        <div class={styles.macroWindows}>
          Только: {props.macros.include_title?.join(", ")}
        </div>
      )}
    </button>
  );
}

export function MacrosList() {
  const navigate = useNavigate();
  const app = useStore($app);

  return (
    <ViewPort
      title="Список макросов"
      subContent={
        <button
          onClick={() => navigate("/macros/create")}
          class={btnStyles.titleBtn}
        >
          Создать макрос
        </button>
      }
    >
      {!app().macros?.length && (
        <div class="h-full flex items-center justify-center font-extrabold text-4xl text-neutral-500">
          Макросов нет, <br />
          Тишина в коде царит, <br />
          Простота везде.
        </div>
      )}

      <div class={styles.list}>
        <For each={app().macros}>
          {(m) => (
            <Macros
              macros={m}
              onToggleDisable={() => toggleMacroDisabled(m.id)}
            />
          )}
        </For>
      </div>
    </ViewPort>
  );
}
