import { Link } from "react-router-dom";
import { useStore } from "@nanostores/react";
import { FaEdit } from "react-icons/fa";
import { FaTrash } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { $app, deleteMacro } from "../stores/app";
import { getKeyNames } from "../constants/keys";
import { MACRO_TYPE_NAMES, MacroType } from "../constants/macro";

export function List() {
  const { macros } = useStore($app);
  const navigate = useNavigate();

  return (
    <div>
      <div className="flex flex-col gap-y-2">
        {macros.length === 0 && <div className="empty">Список пуст</div>}

        {macros.map(macro => (
          <button
            key={macro.id}
            onClick={() => navigate(`/macro/${macro.id}`)}
            className="bg-zinc-900 rounded-xl p-4 flex flex-col gap-y-1 shadow-sm"
          >
            <div className="flex items-center gap-x-2">
              <div className="text-sm truncate font-bold flex-1">
                <div className="truncate min-w-0 flex items-center gap-x-1">
                  <span className="text-xs bg-pink-300 px-1.5 py-0 text-pink-900 rounded-lg mb-0.5">
                    {MACRO_TYPE_NAMES[macro.type as MacroType] || "Неизвестный"}
                  </span>
                  <span className="truncate min-w-0 text-sm">{macro.name}</span>
                </div>
              </div>
            </div>

            <div className="flex flex-wrap items-center gap-1">
              {getKeyNames(macro.activation_keys).map((keyName, index, arr) => (
                <>
                  <div key={index} className="px-1.5 py-0.5 bg-zinc-800 rounded text-xs">
                    {keyName}
                  </div>
                  {index < arr.length - 1 && <div className="text-xs">+</div>}
                </>
              ))}
            </div>

            <div className="text-xs text-zinc-400">
              {macro.include_titles || "Все окна, без ограничения"}
            </div>
          </button>
        ))}
      </div>
    </div>
  );
}
