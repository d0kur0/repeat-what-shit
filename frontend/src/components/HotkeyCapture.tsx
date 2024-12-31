import { useEffect, useState } from "react";
import { EventsOff, EventsOn } from "../../wailsjs/runtime/runtime";
import { StartCapture, StopCapture } from "../../wailsjs/go/main/App";
import { formatKeyCombo } from "../constants/keys";

interface HotkeyCaptureProps {
  value: number[];
  onChange: (keys: number[]) => void;
}

export function HotkeyCapture({ value, onChange }: HotkeyCaptureProps) {
  const handleFocus = () => {
    StartCapture();

    // Блокируем скролл
    const preventWheel = (e: Event) => e.preventDefault();
    document.addEventListener("wheel", preventWheel, { passive: false });

    // Блокируем клик колесом
    const preventMiddleClick = (e: MouseEvent) => {
      if (e.button === 1) {
        e.preventDefault();
      }
    };
    document.addEventListener("mousedown", preventMiddleClick);

    EventsOn("combo_captured", (combo: number[]) => {
      if (combo.some(code => code === 513)) {
        return;
      }
      onChange(combo);
    });

    // Сохраняем функции для последующего удаления
    const cleanup = () => {
      document.removeEventListener("wheel", preventWheel);
      document.removeEventListener("mousedown", preventMiddleClick);
    };

    // Добавляем обработчик blur для очистки
    const input = document.activeElement;
    if (input) {
      const handleBlurOnce = () => {
        cleanup();
        input.removeEventListener("blur", handleBlurOnce);
      };
      input.addEventListener("blur", handleBlurOnce);
    }
  };

  const handleBlur = () => {
    StopCapture();
  };

  useEffect(() => {
    EventsOff("combo_captured");
  }, []);

  return (
    <input
      type="text"
      value={value.length > 0 ? formatKeyCombo(value) : ""}
      readOnly
      onFocus={handleFocus}
      onBlur={handleBlur}
      onClick={() => onChange([513])}
      onContextMenu={e => e.preventDefault()}
      onKeyDown={e => e.preventDefault()}
      onDoubleClick={e => e.preventDefault()}
      onWheel={e => e.preventDefault()}
      className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg focus:outline-none focus:border-violet-500 cursor-pointer select-none"
      placeholder="Нажмите для выбора комбинации..."
    />
  );
}
