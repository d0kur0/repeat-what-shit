export const KEYBOARD_KEYS = {
  SHIFT: 0x10,
  CTRL: 0x11,
  ALT: 0x12,
  "LEFT SHIFT": 0xa0,
  "RIGHT SHIFT": 0xa1,
  "LEFT CTRL": 0xa2,
  "RIGHT CTRL": 0xa3,
  "LEFT ALT": 0xa4,
  "RIGHT ALT": 0xa5,

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

  ESCAPE: 0x1b,
  TAB: 0x09,
  "CAPS LOCK": 0x14,
  "SCROLL LOCK": 0x91,
  "NUM LOCK": 0x90,
  BACKSPACE: 0x08,
  ENTER: 0x0d,
  SPACE: 0x20,
  "PAGE UP": 0x21,
  "PAGE DOWN": 0x22,
  END: 0x23,
  HOME: 0x24,
  INSERT: 0x2d,
  DELETE: 0x2e,

  LEFT: 0x25,
  UP: 0x26,
  RIGHT: 0x27,
  DOWN: 0x28,

  "0": 0x30,
  "1": 0x31,
  "2": 0x32,
  "3": 0x33,
  "4": 0x34,
  "5": 0x35,
  "6": 0x36,
  "7": 0x37,
  "8": 0x38,
  "9": 0x39,

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

  "NUMPAD 0": 0x60,
  "NUMPAD 1": 0x61,
  "NUMPAD 2": 0x62,
  "NUMPAD 3": 0x63,
  "NUMPAD 4": 0x64,
  "NUMPAD 5": 0x65,
  "NUMPAD 6": 0x66,
  "NUMPAD 7": 0x67,
  "NUMPAD 8": 0x68,
  "NUMPAD 9": 0x69,
  MULTIPLY: 0x6a,
  ADD: 0x6b,
  SUBTRACT: 0x6d,
  DECIMAL: 0x6e,
  DIVIDE: 0x6f,

  SEMICOLON: 0xba,
  EQUALS: 0xbb,
  COMMA: 0xbc,
  MINUS: 0xbd,
  PERIOD: 0xbe,
  SLASH: 0xbf,
  BACKQUOTE: 0xc0,
  BRACKET_LEFT: 0xdb,
  BACKSLASH: 0xdc,
  BRACKET_RIGHT: 0xdd,
  QUOTE: 0xde,

  "VOLUME MUTE": 0xad,
  "VOLUME DOWN": 0xae,
  "VOLUME UP": 0xaf,
  "MEDIA NEXT": 0xb0,
  "MEDIA PREV": 0xb1,
  "MEDIA STOP": 0xb2,
  "MEDIA PLAY PAUSE": 0xb3,

  "PRINT SCREEN": 0x2c,
  "PAUSE BREAK": 0x13,
  WIN: 0x5b,
  MENU: 0x5d,

  ЛКМ: 0x0201,
  ПКМ: 0x0204,
  "Колесо мыши": 0x0207,

  XBUTTON1: 0x020b | (0x0001 << 16),
  XBUTTON2: 0x020b | (0x0002 << 16),

  "Колесо вверх": 0x020a | 0x10000,
  "Колесо вниз": 0x020a | 0x20000,
};

export const MODIFIERS = [
  KEYBOARD_KEYS.SHIFT,
  KEYBOARD_KEYS["LEFT SHIFT"],
  KEYBOARD_KEYS["RIGHT SHIFT"],
  KEYBOARD_KEYS.CTRL,
  KEYBOARD_KEYS["LEFT CTRL"],
  KEYBOARD_KEYS["RIGHT CTRL"],
  KEYBOARD_KEYS.ALT,
  KEYBOARD_KEYS["LEFT ALT"],
  KEYBOARD_KEYS["RIGHT ALT"],
];

export function getKeyName(code: number): string {
  const key = Object.keys(KEYBOARD_KEYS).find(
    (key) => KEYBOARD_KEYS[key as keyof typeof KEYBOARD_KEYS] === code
  );

  return key || `Клавиша (${code})`;
}

export function sortKeyCombo(combo: number[]): number[] {
  return [...combo].sort((a, b) => {
    const aIsModifier = MODIFIERS.includes(a);
    const bIsModifier = MODIFIERS.includes(b);

    if (aIsModifier && !bIsModifier) return -1;
    if (!aIsModifier && bIsModifier) return 1;

    return a - b;
  });
}
