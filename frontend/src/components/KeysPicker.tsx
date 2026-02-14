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
  let rootRef!: HTMLButtonElement;

  const suppress = (e: Event) => {
    if (!focused()) return;
    if (e.type === "mousedown" && !rootRef.contains(e.target as Node)) {
      rootRef.blur();
      return;
    }
    e.preventDefault();
    e.stopPropagation();
  };

  const suppressedEvents = [
    "keydown", "keyup", "keypress",
    "mousedown", "mouseup", "auxclick",
    "contextmenu", "dragstart",
  ] as const;

  const handleFocus = () => {
    setFocused(true);
    StartCapture();

    for (const evt of suppressedEvents)
      document.addEventListener(evt, suppress, true);
    document.addEventListener("wheel", suppress, { passive: false, capture: true });

    unsubscribe = Events.On("captured_combo", (ev) => {
      const combo = ev.data as number[];
      if (!focused()) return;
      if (combo.length === 1 && combo[0] === KEYBOARD_KEYS.ЛКМ) return;
      props.onChange(combo);
    });
  };

  const cleanup = () => {
    for (const evt of suppressedEvents)
      document.removeEventListener(evt, suppress, true);
    document.removeEventListener("wheel", suppress, true);
  };

  const handleBlur = () => {
    setFocused(false);
    StopCapture();
    safeUnsubscribe(unsubscribe);
    cleanup();
  };

  const handleClick = () => {
    if (!focused()) return;
    props.onChange([KEYBOARD_KEYS.ЛКМ]);
  };

  onCleanup(() => {
    safeUnsubscribe(unsubscribe);
    cleanup();
  });

  return (
    <button
      ref={rootRef}
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
