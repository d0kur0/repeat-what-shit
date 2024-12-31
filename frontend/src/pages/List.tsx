import { useStore } from "@nanostores/react";
import { FaEdit } from "react-icons/fa";
import { FaTrash } from "react-icons/fa";
import { useNavigate } from "react-router-dom";
import { $app, deleteMacro } from "../stores/app";

export function List() {
  const { macros } = useStore($app);
  const navigate = useNavigate();

  return (
    <div>
      <div className="flex justify-between items-center">
        <div className="font-bold txt-sm">Мои макросы</div>
        <button
          type="button"
          onClick={() => navigate("/macro/create")}
          className="py-1 px-2 font-semibold text-violet-200 inline-flex items-center gap-x-2 text-xs rounded-lg border border-transparent bg-violet-600 hover:bg-violet-700 focus:outline-none disabled:opacity-50 disabled:pointer-events-none"
        >
          Добавить макрос
        </button>
      </div>

      <div className="grid grid-cols-4 gap-4 mt-4 sm:grid-cols-3">
        {macros.length === 0 && <div className="empty">Список пуст</div>}

        {macros.map(macro => (
          <div
            key={macro.id}
            className="bg-zinc-900 rounded-xl p-4 flex flex-col gap-y-3 shadow-sm"
          >
            <div className="flex items-center gap-x-2">
              <div className="text-sm truncate font-bold flex-1">
                <div className="truncate min-w-0">{macro.name}</div>
              </div>
              <div className="flex items-center gap-x-1">
                <button
                  type="button"
                  className="btn"
                  data-hs-tooltip="true"
                  data-hs-tooltip-placement="top"
                  data-hs-tooltip-content="Редактировать"
                >
                  <FaEdit />
                </button>
                <button
                  type="button"
                  className="btn"
                  data-hs-tooltip="true"
                  data-hs-tooltip-placement="top"
                  data-hs-tooltip-content="Удалить"
                  onClick={() => deleteMacro(macro.id)}
                >
                  <FaTrash />
                </button>
              </div>
            </div>

            <div className="flex items-center gap-x-1">
              {macro.activation_keys.map((key, index) => (
                <div key={index} className="kbd">
                  {key}
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
