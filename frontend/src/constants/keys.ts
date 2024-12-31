// Коды клавиш клавиатуры
export const KEYBOARD_KEYS = {
  // Модификаторы
  SHIFT: 0x10,
  CTRL: 0x11,
  ALT: 0x12,
  LEFT_SHIFT: 0xa0,
  RIGHT_SHIFT: 0xa1,
  LEFT_CTRL: 0xa2,
  RIGHT_CTRL: 0xa3,
  LEFT_ALT: 0xa4,
  RIGHT_ALT: 0xa5,

  // Функциональные клавиши
  F1: 0x70,
  F2: 0x71,
  F3: 0x72,
  F4: 0x73,
  F5: 0x74,
  F6: 0x75,
  F7: 0x76,
  F8: 0x77,
  F9: 0x78,
  F10: 0x79,
  F11: 0x7a,
  F12: 0x7b,

  // Управляющие клавиши
  ESCAPE: 0x1b,
  TAB: 0x09,
  CAPS_LOCK: 0x14,
  SCROLL_LOCK: 0x91,
  NUM_LOCK: 0x90,
  BACKSPACE: 0x08,
  ENTER: 0x0d,
  SPACE: 0x20,
  PAGE_UP: 0x21,
  PAGE_DOWN: 0x22,
  END: 0x23,
  HOME: 0x24,
  INSERT: 0x2d,
  DELETE: 0x2e,

  // Стрелки
  LEFT: 0x25,
  UP: 0x26,
  RIGHT: 0x27,
  DOWN: 0x28,

  // Цифры (основная клавиатура)
  KEY_0: 0x30,
  KEY_1: 0x31,
  KEY_2: 0x32,
  KEY_3: 0x33,
  KEY_4: 0x34,
  KEY_5: 0x35,
  KEY_6: 0x36,
  KEY_7: 0x37,
  KEY_8: 0x38,
  KEY_9: 0x39,

  // Буквы
  A: 0x41,
  B: 0x42,
  C: 0x43,
  D: 0x44,
  E: 0x45,
  F: 0x46,
  G: 0x47,
  H: 0x48,
  I: 0x49,
  J: 0x4a,
  K: 0x4b,
  L: 0x4c,
  M: 0x4d,
  N: 0x4e,
  O: 0x4f,
  P: 0x50,
  Q: 0x51,
  R: 0x52,
  S: 0x53,
  T: 0x54,
  U: 0x55,
  V: 0x56,
  W: 0x57,
  X: 0x58,
  Y: 0x59,
  Z: 0x5a,

  // Numpad
  NUMPAD0: 0x60,
  NUMPAD1: 0x61,
  NUMPAD2: 0x62,
  NUMPAD3: 0x63,
  NUMPAD4: 0x64,
  NUMPAD5: 0x65,
  NUMPAD6: 0x66,
  NUMPAD7: 0x67,
  NUMPAD8: 0x68,
  NUMPAD9: 0x69,
  MULTIPLY: 0x6a,
  ADD: 0x6b,
  SUBTRACT: 0x6d,
  DECIMAL: 0x6e,
  DIVIDE: 0x6f,

  // Символы
  SEMICOLON: 0xba, // ;
  EQUALS: 0xbb, // =
  COMMA: 0xbc, // ,
  MINUS: 0xbd, // -
  PERIOD: 0xbe, // .
  SLASH: 0xbf, // /
  BACKQUOTE: 0xc0, // `
  BRACKET_LEFT: 0xdb, // [
  BACKSLASH: 0xdc, // \
  BRACKET_RIGHT: 0xdd, // ]
  QUOTE: 0xde, // '

  // Медиа клавиши
  VOLUME_MUTE: 0xad,
  VOLUME_DOWN: 0xae,
  VOLUME_UP: 0xaf,
  MEDIA_NEXT: 0xb0,
  MEDIA_PREV: 0xb1,
  MEDIA_STOP: 0xb2,
  MEDIA_PLAY_PAUSE: 0xb3,

  // Специальные
  PRINT_SCREEN: 0x2c,
  PAUSE_BREAK: 0x13,
  WIN: 0x5b,
  MENU: 0x5d,
} as const;

// Коды событий мыши (Windows API)
export const MOUSE_EVENTS = {
  // Основные кнопки
  LEFT_MOUSE: 0x0201, // WM_LBUTTONDOWN
  RIGHT_MOUSE: 0x0204, // WM_RBUTTONDOWN
  MIDDLE_MOUSE: 0x0207, // WM_MBUTTONDOWN

  // Дополнительные кнопки
  XBUTTON1: 0x020b | (0x0001 << 16), // WM_XBUTTONDOWN с XBUTTON1
  XBUTTON2: 0x020b | (0x0002 << 16), // WM_XBUTTONDOWN с XBUTTON2

  // Колесо мыши
  WHEEL_UP: 0x020a | 0x10000, // WM_MOUSEWHEEL (прокрутка вверх)
  WHEEL_DOWN: 0x020a | 0x20000, // WM_MOUSEWHEEL (прокрутка вниз)
} as const;

// Группы клавиш для упрощения проверки
const KEY_GROUPS: Record<string, readonly number[]> = {
  SHIFT: [KEYBOARD_KEYS.SHIFT, KEYBOARD_KEYS.LEFT_SHIFT, KEYBOARD_KEYS.RIGHT_SHIFT],
  CTRL: [KEYBOARD_KEYS.CTRL, KEYBOARD_KEYS.LEFT_CTRL, KEYBOARD_KEYS.RIGHT_CTRL],
  ALT: [KEYBOARD_KEYS.ALT, KEYBOARD_KEYS.LEFT_ALT, KEYBOARD_KEYS.RIGHT_ALT],
} as const;

// Функция для получения базового названия модификатора
function getModifierBaseName(code: number): string | null {
  if (KEY_GROUPS.SHIFT.includes(code)) return "SHIFT";
  if (KEY_GROUPS.CTRL.includes(code)) return "CTRL";
  if (KEY_GROUPS.ALT.includes(code)) return "ALT";
  return null;
}

// Функция для получения человекочитаемого названия клавиши
export function getKeyName(code: number): string {
  // Сначала проверяем модификаторы
  const modifierName = getModifierBaseName(code);
  if (modifierName) return modifierName;

  // Поиск в клавиатурных кодах
  for (const [key, value] of Object.entries(KEYBOARD_KEYS)) {
    if (value === code) {
      // Пропускаем модификаторы, так как они уже обработаны выше
      if (key.includes("LEFT_") || key.includes("RIGHT_")) {
        const baseName = key.replace("LEFT_", "").replace("RIGHT_", "");
        if (["SHIFT", "CTRL", "ALT"].includes(baseName)) {
          continue;
        }
      }
      return key.replace(/_/g, " ");
    }
  }

  // Поиск в кодах мыши
  for (const [key, value] of Object.entries(MOUSE_EVENTS)) {
    if (value === code) {
      return key.replace(/_/g, " ");
    }
  }

  return `Unknown (${code})`;
}

// Функция для сортировки комбинации клавиш (модификаторы всегда в начале)
export function sortKeyCombo(combo: number[]): number[] {
  return combo.sort((a, b) => {
    const aIsModifier = getModifierBaseName(a) !== null;
    const bIsModifier = getModifierBaseName(b) !== null;

    if (aIsModifier && !bIsModifier) return -1;
    if (!aIsModifier && bIsModifier) return 1;
    return a - b;
  });
}

// Функция для форматирования комбинации клавиш
export function formatKeyCombo(combo: number[]): string {
  return sortKeyCombo(combo).map(getKeyName).join(" + ");
}

// Функция для получения названий клавиш
export function getKeyNames(combo: number[]): string[] {
  return sortKeyCombo(combo).map(getKeyName);
}
