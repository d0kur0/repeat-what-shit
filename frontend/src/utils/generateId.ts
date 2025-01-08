export function generateId() {
  return +new Date() + Math.random().toString(36).substring(2, 15);
}
