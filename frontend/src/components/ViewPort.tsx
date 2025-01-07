import { JSX } from "solid-js/jsx-runtime";
import styles from "./ViewPort.module.css";
import appIcon from "../assets/appicon.png";
import { FaSolidArrowRightLong, FaSolidXmark } from "solid-icons/fa";
import { createSignal, onMount } from "solid-js";
import { GetVersion } from "../../bindings/repeat-what-shit/internal/app";
import { useNavigate } from "@solidjs/router";
import { Application } from "@wailsio/runtime";

type ViewPortProps = {
  title?: string;
  children: JSX.Element;
  subContent?: JSX.Element;
};

export function ViewPort({ title, children, subContent }: ViewPortProps) {
  const [getVersion, setVersion] = createSignal("0.0.0");
  const navigate = useNavigate();

  onMount(async () => {
    const version = await GetVersion();
    setVersion(version);
  });

  const handleClose = () => {
    Application.Hide();
  };

  return (
    <div class={styles.root}>
      <div class={styles.title}>
        <button class={styles.app} onClick={() => navigate("/")}>
          <div class={styles.appIcon}>
            <img src={appIcon} alt="Macro Editor" />
          </div>
          <div class={styles.appName}>
            repeat that shit <span class={styles.version}>v{getVersion()}</span>
          </div>
          <FaSolidArrowRightLong class={styles.appArrow} />
          <div class={styles.appSubName}>{title}</div>
        </button>

        <div class="flex-1" />

        <div class={styles.subContent}>{subContent}</div>

        <button onClick={handleClose} class={styles.closeButton}>
          <FaSolidXmark />
        </button>
      </div>
      <div class={styles.content}>{children}</div>
    </div>
  );
}
