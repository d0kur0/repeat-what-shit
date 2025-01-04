export const KEYBOARD_KEYS = {
  SHIFT: 0x10,
  CTRL: 0x11,
  ALT: 0x12,
  LEFT_SHIFT: 0xa0,
  RIGHT_SHIFT: 0xa1,
  LEFT_CTRL: 0xa2,
  RIGHT_CTRL: 0xa3,
  LEFT_ALT: 0xa4,
  RIGHT_ALT: 0xa5,

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

  LEFT: 0x25,
  UP: 0x26,
  RIGHT: 0x27,
  DOWN: 0x28,

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

  VOLUME_MUTE: 0xad,
  VOLUME_DOWN: 0xae,
  VOLUME_UP: 0xaf,
  MEDIA_NEXT: 0xb0,
  MEDIA_PREV: 0xb1,
  MEDIA_STOP: 0xb2,
  MEDIA_PLAY_PAUSE: 0xb3,

  PRINT_SCREEN: 0x2c,
  PAUSE_BREAK: 0x13,
  WIN: 0x5b,
  MENU: 0x5d,

  LEFT_MOUSE: 0x0201,
  RIGHT_MOUSE: 0x0204,
  MIDDLE_MOUSE: 0x0207,

  XBUTTON1: 0x020b | (0x0001 << 16),
  XBUTTON2: 0x020b | (0x0002 << 16),

  WHEEL_UP: 0x020a | 0x10000,
  WHEEL_DOWN: 0x020a | 0x20000,
};

export const MODIFIERS = [
  KEYBOARD_KEYS.SHIFT,
  KEYBOARD_KEYS.LEFT_SHIFT,
  KEYBOARD_KEYS.RIGHT_SHIFT,
  KEYBOARD_KEYS.CTRL,
  KEYBOARD_KEYS.LEFT_CTRL,
  KEYBOARD_KEYS.RIGHT_CTRL,
  KEYBOARD_KEYS.ALT,
  KEYBOARD_KEYS.LEFT_ALT,
  KEYBOARD_KEYS.RIGHT_ALT,
];

export function getKeyName(code: number): string {
  const key = Object.keys(KEYBOARD_KEYS).find(
    (key) => KEYBOARD_KEYS[key as keyof typeof KEYBOARD_KEYS] === code
  );

  return key || `Unknown (${code})`;
}

export function sortKeyCombo(combo: number[]): number[] {
  return combo.sort((a, b) => {
    const aIsModifier = MODIFIERS.includes(a);
    const bIsModifier = MODIFIERS.includes(b);

    if (aIsModifier && !bIsModifier) return -1;
    if (!aIsModifier && bIsModifier) return 1;

    return a - b;
  });
}
