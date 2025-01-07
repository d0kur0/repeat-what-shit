import { createSignal, For, onCleanup, onMount } from "solid-js";
import styles from "./KeysPicker.module.css";
import { getKeyName, KEYBOARD_KEYS, sortKeyCombo } from "../utils/keys";
import { Events } from "@wailsio/runtime";
import {
  StartCapture,
  StopCapture,
} from "../../bindings/repeat-what-shit/internal/app";

type KeysPickerProps = {
  value: number[];
  onChange: (combo: number[]) => void;
};

// Currently, @wailsio/runtime have bug
// if u try unsub for event, what already unsub:
// Cannot read properties of undefined (reading 'filter')
const safeUnsubscribe = (f: Function | undefined) => {
  try {
    f?.();
  } catch {}
};

export function KeysPicker(props: KeysPickerProps) {
  const [focused, setFocused] = createSignal(false);
  let unsubscribe: Function;

  const preventDefaultHandler = (e: Event) => {
    if (focused()) {
      e.preventDefault();
      return false;
    }
  };

  const handleKeyDown = (e: KeyboardEvent) => {
    if (focused() && e.key === "Tab") {
      e.preventDefault();
      return false;
    }
  };

  const handleFocus = () => {
    setFocused(true);
    StartCapture();

    document.addEventListener("wheel", preventDefaultHandler, {
      passive: false,
    });
    document.addEventListener("keydown", handleKeyDown);
    document.addEventListener("contextmenu", preventDefaultHandler);

    unsubscribe = Events.On(
      "captured_combo",
      ({ data: [combo] }: { data: [number[]] }) => {
        if (!focused()) return;
        if (combo.length === 1 && combo[0] === KEYBOARD_KEYS.ЛКМ) return;
        props.onChange(combo);
      }
    );
  };

  const handleBlur = () => {
    setFocused(false);
    StopCapture();
    safeUnsubscribe(unsubscribe);

    document.removeEventListener("wheel", preventDefaultHandler);
    document.removeEventListener("keydown", handleKeyDown);
    document.removeEventListener("contextmenu", preventDefaultHandler);
  };

  const handleClick = () => {
    if (!focused()) return;
    props.onChange([KEYBOARD_KEYS.ЛКМ]);
  };

  onCleanup(() => {
    safeUnsubscribe(unsubscribe);
    document.removeEventListener("wheel", preventDefaultHandler);
    document.removeEventListener("keydown", handleKeyDown);
    document.removeEventListener("contextmenu", preventDefaultHandler);
  });

  return (
    <button
      class={styles.root}
      onBlur={handleBlur}
      onClick={handleClick}
      onFocus={handleFocus}
    >
      {focused() && (
        <div class={styles.label}>
          Для остановки захвата просто кликни вне данного поля
        </div>
      )}
      {!!props.value.length || "Нажми для захвата клавиш"}
      {!!props.value.length && (
        <div class={styles.keys}>
          <For each={sortKeyCombo(props.value)}>
            {(key, idx) => (
              <>
                <div class={styles.key}>{getKeyName(key)}</div>
                {idx() !== props.value.length - 1 && (
                  <div class={styles.separator}>+</div>
                )}
              </>
            )}
          </For>
        </div>
      )}
    </button>
  );
}
