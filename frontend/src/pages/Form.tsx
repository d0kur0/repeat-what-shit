import { useNavigate } from "react-router-dom";
import { FaArrowLeft, FaTrash } from "react-icons/fa";
import { useState } from "react";
import { formatKeyCombo } from "../constants/keys";
import { HotkeyCapture } from "../components/HotkeyCapture";

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
  const [actions, setActions] = useState<Array<{ keys: number[]; delay: number }>>([]);

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
            <HotkeyCapture value={capturedKeys} onChange={setCapturedKeys} />
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

          {/* Действия */}
          <div>
            <div className="flex items-center justify-between mb-3">
              <label className="block text-sm font-medium">Действия</label>
              <button
                onClick={() => setActions([...actions, { keys: [], delay: 0 }])}
                className="px-3 py-1 bg-violet-600 hover:bg-violet-700 rounded-lg text-sm transition-colors"
              >
                Добавить действие
              </button>
            </div>

            {actions.length === 0 ? (
              <div className="text-sm text-zinc-400 text-center py-8 bg-zinc-800/50 rounded-lg">
                Пока нет добавленных действий
              </div>
            ) : (
              <div className="space-y-2">
                {actions.map((action, index) => (
                  <div key={index} className="flex items-center gap-x-4 bg-zinc-800 rounded-lg p-3">
                    <div className="flex-1 flex items-center gap-x-4">
                      <div className="flex-1">
                        <HotkeyCapture
                          value={action.keys}
                          onChange={keys => {
                            const newActions = [...actions];
                            newActions[index].keys = keys;
                            setActions(newActions);
                          }}
                        />
                      </div>
                      <input
                        type="number"
                        min="0"
                        value={action.delay}
                        onChange={e => {
                          const newActions = [...actions];
                          newActions[index].delay = Math.max(0, parseInt(e.target.value) || 0);
                          setActions(newActions);
                        }}
                        className="w-24 px-3 py-2 bg-zinc-700 border border-zinc-600 rounded-lg focus:outline-none focus:border-violet-500 text-sm"
                        placeholder="Задержка (мс)"
                      />
                    </div>
                    <button
                      onClick={() => setActions(actions?.filter((_, i) => i !== index))}
                      className="p-2 hover:bg-zinc-700 rounded-lg transition-colors text-red-500"
                    >
                      <FaTrash />
                    </button>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
