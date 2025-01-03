import { useNavigate, useParams } from "react-router-dom";
import { FaArrowLeft, FaTrash } from "react-icons/fa";
import { useState, useMemo, useEffect } from "react";
import { formatKeyCombo } from "../constants/keys";
import { HotkeyCapture } from "../components/HotkeyCapture";
import { main } from "../../wailsjs/go/models";
import { addMacro, updateMacro, $app } from "../stores/app";
import { v4 as uuidv4 } from "uuid";
import { useStore } from "@nanostores/react";
import { MacroType, MACRO_TYPE_NAMES, DEFAULT_MACRO } from "../constants/macro";

const MACRO_TYPE_DESCRIPTIONS: Record<MacroType, string> = {
  [MacroType.Sequence]: "Срабатывает при каждом нажатии комбинации клавиш",
  [MacroType.Toggle]: "Включается и выключается при нажатии комбинации клавиш",
  [MacroType.Hold]: "Выполняется многократно, пока зажата комбинация клавиш",
};

export function Form() {
  const navigate = useNavigate();
  const { id } = useParams();
  const { macros } = useStore($app);
  const [macro, setMacro] = useState<Omit<main.Macro, "id">>(DEFAULT_MACRO);

  useEffect(() => {
    if (id) {
      const existingMacro = macros.find(m => m.id === id);
      if (existingMacro) {
        const { id: _, ...macroWithoutId } = existingMacro;
        setMacro(macroWithoutId);
      }
    }
  }, [id, macros]);

  const isValid = useMemo(() => {
    return (
      macro.name.trim() !== "" && // Имя не пустое
      macro.activation_keys.length > 0 && // Есть клавиши активации
      macro.actions.length > 0 && // Есть хотя бы одно действие
      macro.actions.every(action => action.keys.length > 0) // У каждого действия есть клавиши
    );
  }, [macro]);

  const setMacroField = (updates: Partial<Omit<main.Macro, "id">>) => {
    setMacro(current => ({ ...current, ...updates }));
  };

  const updateAction = (index: number, updates: Partial<main.MacroAction>) => {
    setMacro(current => ({
      ...current,
      actions: current.actions.map((action, i) =>
        i === index ? { ...action, ...updates } : action
      ),
    }));
  };

  const addAction = () => {
    setMacro(current => ({
      ...current,
      actions: [...current.actions, { keys: [], delay: 0 }],
    }));
  };

  const removeAction = (index: number) => {
    setMacro(current => ({
      ...current,
      actions: current.actions.filter((_, i) => i !== index),
    }));
  };

  const handleSubmit = async () => {
    if (id) {
      await updateMacro({ ...macro, id });
    } else {
      await addMacro({ ...macro, id: uuidv4() });
    }
    navigate("/");
  };

  return (
    <div>
      <div className="">
        <div className="space-y-6">
          {/* Имя макроса */}
          <div>
            <label htmlFor="name" className="block text-sm font-medium mb-2">
              Название макроса
            </label>
            <input
              type="text"
              id="name"
              value={macro.name}
              onChange={e => setMacroField({ name: e.target.value })}
              className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg focus:outline-none focus:border-violet-500"
              placeholder="Например: Быстрый прыжок"
            />
            {macro.name.trim() === "" && (
              <div className="mt-1.5 text-xs text-red-400">Введите название макроса</div>
            )}
          </div>

          {/* Комбинация клавиш */}
          <div>
            <label htmlFor="hotkey" className="block text-sm font-medium mb-2">
              Комбинация клавиш
            </label>
            <HotkeyCapture
              value={macro.activation_keys}
              onChange={activation_keys => setMacroField({ activation_keys })}
            />
            {macro.activation_keys.length === 0 && (
              <div className="mt-1.5 text-xs text-red-400">
                Выберите комбинацию клавиш для активации
              </div>
            )}
          </div>

          {/* Тип макроса */}
          <div className="flex flex-col gap-4">
            <div>
              <label className="block text-sm font-medium mb-3">Тип макроса</label>
              <div className="space-y-3">
                {Object.entries(MACRO_TYPE_DESCRIPTIONS).map(([type, description]) => (
                  <div key={type} className="flex">
                    <div className="flex h-7">
                      <input
                        id={type}
                        name="macro_type"
                        type="radio"
                        checked={macro.type === parseInt(type)}
                        onChange={() => setMacroField({ type: parseInt(type) as MacroType })}
                        className="h-5 w-5 border-zinc-700 accent-violet-500 bg-zinc-800 cursor-pointer"
                      />
                    </div>
                    <div
                      className="ml-3 cursor-pointer"
                      onClick={() => setMacroField({ type: parseInt(type) as MacroType })}
                    >
                      <label htmlFor={type} className="block text-sm font-medium">
                        {MACRO_TYPE_NAMES[parseInt(type) as MacroType]}
                      </label>
                      <p className="text-zinc-400 text-sm">{description}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Привязка к окнам */}
            <div>
              <label htmlFor="include_titles" className="block text-sm font-medium mb-2">
                Привязка к окнам
              </label>
              <input
                type="text"
                id="include_titles"
                value={macro.include_titles}
                onChange={e => setMacroField({ include_titles: e.target.value })}
                className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg focus:outline-none focus:border-violet-500"
                placeholder="Например: Notepad++, Chrome, code.exe"
              />
              <div className="mt-1 text-xs text-zinc-400">
                Укажите названия окон или их классы через запятую. Оставьте пустым для работы во
                всех окнах
              </div>
            </div>

            {/* Действия */}
            <div>
              <div className="flex items-center mb-6 gap-x-2">
                <label className="block text-sm font-medium">Действия</label>
                <button
                  onClick={addAction}
                  className="px-2 py-0.5 bg-violet-600 hover:bg-violet-700 rounded-lg text-xs transition-colors"
                >
                  Добавить
                </button>
              </div>

              {macro.actions.length === 0 ? (
                <div className="border-red-400 border-dashed border-2 rounded-lg p-6 text-center">
                  <div className="text-sm text-red-400">Добавьте хотя бы одно действие</div>
                </div>
              ) : (
                <div className="space-y-5">
                  {macro.actions.map((action, index) => (
                    <div key={index} className="flex flex-col gap-y-1.5">
                      <div className="flex items-center gap-x-4">
                        <div className="flex-1 flex gap-x-4">
                          <div className="flex-1">
                            <HotkeyCapture
                              value={action.keys}
                              onChange={keys => updateAction(index, { keys })}
                            />
                          </div>
                          <div className="relative py-1.5 px-3 w-42 bg-zinc-800 border border-zinc-700 rounded-lg">
                            <span className="absolute -top-2 left-2 px-1 text-[10px] leading-none text-zinc-400 bg-zinc-800">
                              Пауза после (ms)
                            </span>
                            <div
                              className="w-full flex justify-between items-center gap-x-2"
                              data-hs-input-number
                            >
                              <div>
                                <input
                                  className={`p-0 bg-transparent border-0 text-white text-sm outline-none [&::-webkit-inner-spin-button]:appearance-none [&::-webkit-outer-spin-button]:appearance-none w-full`}
                                  style={{ WebkitAppearance: "none", MozAppearance: "textfield" }}
                                  type="number"
                                  value={action.delay}
                                  onChange={e =>
                                    updateAction(index, {
                                      delay: +e.target.value,
                                    })
                                  }
                                  data-hs-input-number-input
                                />
                              </div>
                              <div className="flex justify-end items-center gap-x-1">
                                <button
                                  type="button"
                                  onClick={() =>
                                    updateAction(index, {
                                      delay: Math.max(0, action.delay - 100),
                                    })
                                  }
                                  className="size-5 inline-flex justify-center items-center text-sm font-medium rounded-full border border-zinc-600 bg-zinc-700 text-white hover:bg-zinc-600 focus:outline-none disabled:opacity-50 disabled:pointer-events-none"
                                  tabIndex={-1}
                                  aria-label="Decrease"
                                  data-hs-input-number-decrement
                                >
                                  <svg
                                    className="size-3"
                                    xmlns="http://www.w3.org/2000/svg"
                                    width="24"
                                    height="24"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                  >
                                    <path d="M5 12h14"></path>
                                  </svg>
                                </button>
                                <button
                                  type="button"
                                  onClick={() =>
                                    updateAction(index, {
                                      delay: action.delay + 100,
                                    })
                                  }
                                  className="size-5 inline-flex justify-center items-center text-sm font-medium rounded-full border border-zinc-600 bg-zinc-700 text-white hover:bg-zinc-600 focus:outline-none disabled:opacity-50 disabled:pointer-events-none"
                                  tabIndex={-1}
                                  aria-label="Increase"
                                  data-hs-input-number-increment
                                >
                                  <svg
                                    className="size-3"
                                    xmlns="http://www.w3.org/2000/svg"
                                    width="24"
                                    height="24"
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    strokeWidth="2"
                                    strokeLinecap="round"
                                    strokeLinejoin="round"
                                  >
                                    <path d="M5 12h14"></path>
                                    <path d="M12 5v14"></path>
                                  </svg>
                                </button>
                              </div>
                            </div>
                          </div>
                        </div>

                        <div>
                          <button
                            onClick={() => removeAction(index)}
                            className="p-1.5 text-xs border-red-400 border-2 rounded-lg hover:opacity-75 text-red-300"
                          >
                            <FaTrash />
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>

            {/* Кнопка сохранения */}
            <div className="flex justify-end">
              <button
                onClick={handleSubmit}
                disabled={!isValid}
                className="px-4 py-2 bg-violet-600 hover:bg-violet-700 rounded-lg transition-colors disabled:opacity-50 disabled:pointer-events-none"
              >
                Сохранить
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
