import { ErrorBoundary } from "solid-js";
import { hashIntegration, Router } from "@solidjs/router";
import { Route } from "@solidjs/router";
import { Routes } from "@solidjs/router";
import { MacrosList } from "./pages/MacrosList";
import { MacrosForm } from "./pages/MacrosForm";
import { Settings } from "./pages/Settings";

function ErrorScreen(props: { err: Error }) {
  return (
    <div class="prose max-w-full w-full">
      <h1>Runtime error</h1>

      <h3>Ebat, why?</h3>

      <h4>{props.err.message}</h4>

      <pre class="w-full mt-4">{props.err.stack}</pre>
    </div>
  );
}

export function App() {
  return (
    <ErrorBoundary fallback={(err) => <ErrorScreen err={err as Error} />}>
      <Router source={hashIntegration()}>
        <Routes>
          <Route path="/" component={MacrosList} />
          <Route path="/macros/create" component={MacrosForm} />
          <Route path="/macros/edit/:id" component={MacrosForm} />
          <Route path="/settings" component={Settings} />
        </Routes>
      </Router>
    </ErrorBoundary>
  );
}
