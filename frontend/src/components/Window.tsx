import { FaMinus } from "react-icons/fa6";
import { IoClose } from "react-icons/io5";
import { FiMinimize2 } from "react-icons/fi";

import "./Window.css";

type WindowProps = {
  children: React.ReactNode;
};

export function Window({ children }: WindowProps) {
  return (
    <div className="window">
      <div className="window-title">
        <div className="window-title-text andika" data-draggable>
          repeat what shit
        </div>

        <div className="window-title-buttons">
          <button className="minimize">
            <FaMinus />
          </button>
          <button className="maximize">
            <FiMinimize2 />
          </button>
          <button className="close">
            <IoClose />
          </button>
        </div>
      </div>

      <div className="window-content">{children}</div>
    </div>
  );
}
