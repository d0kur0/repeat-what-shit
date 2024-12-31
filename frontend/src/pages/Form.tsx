import { useNavigate } from "react-router-dom";
import { FaArrowLeft } from "react-icons/fa";
import { useState, useEffect } from "react";
import { EventsOn } from "../../wailsjs/runtime/runtime";
import { StartCapture, StopCapture } from "../../wailsjs/go/main/App";
import { getKeyName, formatKeyCombo } from "../constants/keys";

type MacroType = "sequence" | "whilePressed" | "toggle";

const MACRO_TYPE_DESCRIPTIONS: Record<MacroType, string> = {
  sequence: "Выполнить действия один раз по нажатию клавиши",
  whilePressed: "Повторять действия, пока зажата клавиша",
  toggle: "Повторять действия до повторного нажатия клавиши",
};

const MACRO_TYPE_NAMES: Record<MacroType, string> = {
  sequence: "По нажатию",
  whilePressed: "При удержании",
  toggle: "Переключение",
};

export function Form() {
  const navigate = useNavigate();
  const [name, setName] = useState("");
  const [type, setType] = useState<MacroType>("sequence");
  const [capturedKeys, setCapturedKeys] = useState<number[]>([]);

  useEffect(() => {
    const unsubscribe = EventsOn("combo_captured", (combo: number[]) => {
      console.log("combo_captured", combo);
      setCapturedKeys(combo);
    });

    return () => unsubscribe();
  }, []);

  return (
    <div>
      <div className="flex items-center gap-x-4 mb-6">
        <button
          onClick={() => navigate("/")}
          className="p-2 hover:bg-zinc-800 rounded-lg transition-colors"
        >
          <FaArrowLeft />
        </button>
        <div className="font-bold">Создание макроса</div>
      </div>

      <div className="bg-zinc-900 rounded-xl p-4">
        <div className="space-y-6">
          {/* Имя макроса */}
          <div>
            <label htmlFor="name" className="block text-sm font-medium mb-2">
              Название макроса
            </label>
            <input
              type="text"
              id="name"
              value={name}
              onChange={e => setName(e.target.value)}
              className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg focus:outline-none focus:border-violet-500"
              placeholder="Например: Быстрый прыжок"
            />
          </div>

          {/* Комбинация клавиш */}
          <div>
            <label htmlFor="hotkey" className="block text-sm font-medium mb-2">
              Комбинация клавиш
            </label>
            <input
              type="text"
              id="hotkey"
              value={formatKeyCombo(capturedKeys)}
              readOnly
              onFocus={() => StartCapture()}
              onBlur={() => StopCapture()}
              onKeyDown={e => e.preventDefault()}
              onContextMenu={e => e.preventDefault()}
              onDoubleClick={e => e.preventDefault()}
              className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg focus:outline-none focus:border-violet-500 cursor-pointer select-none"
              placeholder="Нажмите для захвата клавиш..."
            />
          </div>

          {/* Тип макроса */}
          <div>
            <label className="block text-sm font-medium mb-3">Тип макроса</label>
            <div className="space-y-3">
              {(Object.keys(MACRO_TYPE_DESCRIPTIONS) as MacroType[]).map(macroType => (
                <div key={macroType} className="flex">
                  <div className="flex h-7">
                    <input
                      id={macroType}
                      name="macro_type"
                      type="radio"
                      checked={type === macroType}
                      onChange={() => setType(macroType)}
                      className="h-5 w-5 border-zinc-700 accent-violet-500 bg-zinc-800 cursor-pointer"
                    />
                  </div>
                  <div className="ml-3 cursor-pointer" onClick={() => setType(macroType)}>
                    <label htmlFor={macroType} className="block text-sm font-medium">
                      {MACRO_TYPE_NAMES[macroType]}
                    </label>
                    <p className="text-zinc-400 text-sm">{MACRO_TYPE_DESCRIPTIONS[macroType]}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
