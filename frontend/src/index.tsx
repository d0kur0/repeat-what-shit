/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import { App } from "./App";

const root = document.getElementById("root");

import "@wailsio/runtime";

for (const evt of ["mousedown", "mouseup", "auxclick"] as const) {
  document.addEventListener(
    evt,
    (e) => {
      if (e.button === 3 || e.button === 4) {
        e.preventDefault();
        e.stopPropagation();
      }
    },
    true,
  );
}

if (import.meta.env.DEV && !(root instanceof HTMLElement)) {
  throw new Error(
    "Root element not found. Did you forget to add it to your index.html? Or maybe the id attribute got misspelled?",
  );
}

render(() => <App />, root!);
