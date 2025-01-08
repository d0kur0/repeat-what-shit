import { useNavigate, useParams } from "@solidjs/router";
import { ViewPort } from "../components/ViewPort";

import styles from "./MacrosForm.module.css";
import btnStyles from "../styles/Buttons.module.css";
import inputStyles from "../styles/Inputs.module.css";
import modalStyles from "../styles/Modal.module.css";

import { createStore } from "solid-js/store";
import {
  Macro,
  MacroType,
} from "../../bindings/repeat-what-shit/internal/types";
import { createSignal, For, onMount } from "solid-js";
import { KeysPicker } from "../components/KeysPicker";
import { Key } from "@solid-primitives/keyed";
import { generateId } from "../utils/generateId";
import { useStore } from "@nanostores/solid";
import { $app, addMacro, deleteMacro, updateMacro } from "../stores/app";
import { Dialog } from "@kobalte/core/dialog";
import { IoClose, IoCloseSharp } from "solid-icons/io";
import { WindowInfo } from "../../bindings/repeat-what-shit/internal/utils";
import { GetWindowList } from "../../bindings/repeat-what-shit/internal/app";

import unknownFile from "../assets/unknown-file-types.png";

type ValidationError = {
  path: string;
  message: string;
};

const validateMacro = (macro: Macro) => {
  const errors: ValidationError[] = [];

  macro.name ||
    errors.push({ path: "name", message: "Введите название макроса" });

  if (!macro.activation_keys?.length) {
    errors.push({
      path: "activation_keys",
      message: "Выберите хотя бы одну клавишу активации",
    });
  }

  if (!macro.actions || macro.actions.length === 0) {
    errors.push({ path: "actions", message: "Добавьте хотя бы одно действие" });
  } else {
    macro.actions.forEach((action, index) => {
      if (!action.keys || action.keys.length === 0) {
        errors.push({
          path: `actions.${action.id}.keys`,
          message: "Выберите клавиши для действия",
        });
      }
      if (action.delay < 0) {
        errors.push({
          path: `actions.${action.id}.delay`,
          message: "Задержка не может быть отрицательной",
        });
      }
    });
  }

  return errors;
};

export function MacrosForm() {
  const { id } = useParams();
  const isEdit = !!id;

  const [windows, setWindows] = createSignal<WindowInfo[]>([]);

  const app = useStore($app);
  const navigate = useNavigate();

  const [validationErrors, setValidationErrors] = createSignal<
    ValidationError[]
  >([]);

  const getErrorByPath = (path: string) => {
    return validationErrors().find((error) => error.path === path)?.message;
  };

  const [macros, setMacros] = createStore<Macro>({
    id: generateId(),
    name: "",
    type: MacroType.MacroTypeSequence,
    disabled: false,
    actions: [
      {
        id: generateId(),
        keys: [],
        delay: 0,
      },
    ],
    include_title: [],
    activation_keys: [],
  });

  const fetchWindowsList = async () => {
    const windows = await GetWindowList();
    setWindows(windows || []);
  };

  onMount(async () => {
    const macro = app().macros?.find((v) => v.id === id);
    macro && setMacros(macro);
    await fetchWindowsList();
  });

  const updateWindowList = async () => {
    await fetchWindowsList();
  };

  const handleSave = () => {
    setValidationErrors([]);
    const validationErrors = validateMacro(macros);

    if (validationErrors.length) {
      return setValidationErrors(validationErrors);
    }

    isEdit ? updateMacro(macros) : addMacro(macros);
    navigate("/");
  };

  const handleDelete = () => {
    const really = confirm("Точно удалить?");
    if (!really) return;

    deleteMacro(macros.id);
    navigate("/");
  };

  const handleAddIncludeTitle = (title: string) => {
    setMacros("include_title", (v) => [...(v || []), title]);
  };

  return (
    <ViewPort
      subContent={
        <div class="flex gap-3">
          {isEdit && (
            <button
              onClick={handleDelete}
              classList={{
                [btnStyles.titleBtn]: true,
                [btnStyles.titleBtnRemove]: true,
              }}
            >
              Удалить
            </button>
          )}

          <button
            onClick={handleSave}
            classList={{
              [btnStyles.titleBtn]: true,
              [btnStyles.titleBtnSave]: true,
            }}
          >
            Сохранить
          </button>
        </div>
      }
      title={isEdit ? "Редактирование макроса" : "Создание макроса"}
    >
      <div class={inputStyles.inputContainer}>
        <div class={inputStyles.label}>Название макроса</div>
        <input
          value={macros.name}
          class={inputStyles.input}
          placeholder="write something"
          onInput={(e) => setMacros("name", e.target.value)}
        />
        <div class={inputStyles.error}>{getErrorByPath("name")}</div>
      </div>

      <div class={inputStyles.inputContainer}>
        <div class={inputStyles.label}>Клавиши активации</div>
        <KeysPicker
          value={macros.activation_keys || []}
          onChange={(combo) => setMacros("activation_keys", combo)}
        />
        <div class={inputStyles.error}>{getErrorByPath("activation_keys")}</div>
      </div>

      <div class="flex items-center justify-between">
        <div class={styles.includeTitleTitle}>Привязать к окнам</div>

        <Dialog>
          <Dialog.Trigger
            onClick={updateWindowList}
            class={styles.includeTitlesBtn}
          >
            Добавить окно
          </Dialog.Trigger>
          <Dialog.Portal>
            <Dialog.Overlay class={modalStyles.dialogOverlay} />
            <div class={modalStyles.dialogPositioner}>
              <Dialog.Content class={modalStyles.dialogContent}>
                <div class={modalStyles.dialogHeader}>
                  <Dialog.Title class={modalStyles.dialogTitle}>
                    Добавление окна для макроса
                  </Dialog.Title>
                  <Dialog.CloseButton class={modalStyles.dialogCloseButton}>
                    <IoClose />
                  </Dialog.CloseButton>
                </div>
                <Dialog.Description class={modalStyles.dialogDescription}>
                  <div class={styles.windows}>
                    <For each={windows()}>
                      {(v) => {
                        const alreadyAdded = macros.include_title?.includes(
                          v.process
                        );

                        return (
                          <Dialog.CloseButton
                            onClick={() => handleAddIncludeTitle(v.process)}
                            class={styles.window}
                            disabled={alreadyAdded}
                          >
                            <div class={styles.windowIconContainer}>
                              <img
                                class={styles.windowIcon}
                                src={
                                  v.iconBase64
                                    ? `data:image/png;base64,${v.iconBase64}`
                                    : unknownFile
                                }
                              />
                            </div>
                            <div class="flex items-center justify-between flex-1">
                              <div>{v.process}</div>
                              {alreadyAdded && (
                                <div class={styles.alreadyAddedBadge}>
                                  Уже добавлен
                                </div>
                              )}
                            </div>
                          </Dialog.CloseButton>
                        );
                      }}
                    </For>
                  </div>
                </Dialog.Description>
              </Dialog.Content>
            </div>
          </Dialog.Portal>
        </Dialog>
      </div>
      <div class={styles.includeTitleDescription}>
        Возможность указать окна, в которых будет работать макрос
      </div>

      {!!macros.include_title?.length || (
        <div class={styles.includeTitles}>Ничего не выбрано</div>
      )}

      {!!macros.include_title?.length && (
        <div class={styles.includeTitles}>
          <For each={macros.include_title}>
            {(t, id) => (
              <button
                onClick={() => {
                  setMacros(
                    "include_title",
                    (v) => v?.filter((_, idx) => idx !== id()) || []
                  );
                }}
                class={styles.windowBadge}
              >
                {t}
                <IoCloseSharp />
              </button>
            )}
          </For>
        </div>
      )}

      <div class={styles.typesTitle}>Выберите тип макроса</div>

      <div class={styles.types}>
        <button
          onClick={() => setMacros("type", MacroType.MacroTypeSequence)}
          classList={{
            [styles.type]: true,
            [styles.typeActive]: macros.type === MacroType.MacroTypeSequence,
          }}
        >
          <div class={styles.typeTitle}>Один раз по нажатию</div>
          <div class={styles.typeDescription}>
            Выполняется на одно нажатие, повторяется если зажать и подождать
            системного повтора
          </div>
        </button>

        <button
          onClick={() => setMacros("type", MacroType.MacroTypeToggle)}
          classList={{
            [styles.type]: true,
            [styles.typeActive]: macros.type === MacroType.MacroTypeToggle,
          }}
        >
          <div class={styles.typeTitle}>Переключение</div>
          <div class={styles.typeDescription}>
            Первым нажатием активируется циклическое повторение, вторым
            выключается
          </div>
        </button>

        <button
          onClick={() => setMacros("type", MacroType.MacroTypeHold)}
          classList={{
            [styles.type]: true,
            [styles.typeActive]: macros.type === MacroType.MacroTypeHold,
          }}
        >
          <div class={styles.typeTitle}>Зажатие</div>
          <div class={styles.typeDescription}>
            Выполняется пока зажаты клавиши активации, без задержек
          </div>
        </button>
      </div>

      <div class={styles.actionsTitle}>
        <div>Шаги макроса</div>
        <button
          onClick={() =>
            setMacros("actions", (v) => [
              ...(v || []),
              { id: generateId(), keys: [], delay: 0 },
            ])
          }
          class={styles.addActionBtn}
        >
          Добавить шаг
        </button>
      </div>

      <div class={styles.actions}>
        {!!macros.actions?.length && (
          <div
            classList={{
              [styles.action]: true,
              "mb-2 text-xs text-neutral-400": true,
            }}
          >
            <div>Клавиши</div>
            <div>Задержка (ms)</div>
            <div></div>
          </div>
        )}

        <Key
          fallback={
            <div class={styles.actionsFallback}>Добавьте хотя бы один шаг</div>
          }
          each={macros.actions}
          by={(action) => action.id}
        >
          {(action) => (
            <div class={styles.action}>
              <div class={inputStyles.inputContainer}>
                <KeysPicker
                  value={action().keys || []}
                  onChange={(combo) => {
                    setMacros("actions", (v) => {
                      const newActions = [...(v || [])].map((a) =>
                        a.id === action().id ? { ...a, keys: combo } : a
                      );
                      return newActions;
                    });
                  }}
                />
                <div class={inputStyles.error}>
                  {getErrorByPath(`actions.${action().id}.keys`)}
                </div>
              </div>

              <div class={inputStyles.inputContainer}>
                <input
                  type="number"
                  value={action().delay}
                  class={inputStyles.input}
                  onInput={(e) =>
                    setMacros("actions", (v) => {
                      const newActions = [...(v || [])].map((a) =>
                        a.id === action().id
                          ? { ...a, delay: +e.target.value }
                          : a
                      );
                      return newActions;
                    })
                  }
                />
                <div class={inputStyles.error} />
              </div>

              <button
                class={styles.removeActionBtn}
                onClick={() => {
                  setMacros("actions", (v) => {
                    const newActions = [...(v || [])]?.filter(
                      (a) => a.id !== action().id
                    );
                    return newActions;
                  });
                }}
              >
                Удалить
              </button>
            </div>
          )}
        </Key>
      </div>
    </ViewPort>
  );
}
