import { IoAdd, IoArrowBack, IoArrowForward, IoClose, IoTrash } from "react-icons/io5";
import { WindowHide } from "../../wailsjs/runtime/runtime";
import logo from "../assets/appicon.png";
import { deleteMacro } from "../stores/app";

import "./Window.css";
import { useLocation, useNavigate } from "react-router-dom";

type WindowProps = {
  children: React.ReactNode;
};

export function Window({ children }: WindowProps) {
  const location = useLocation();
  const navigate = useNavigate();

  const getPageTitle = () => {
    if (location.pathname === "/") return "Список макросов";
    if (location.pathname === "/macro/create") return "Создание макроса";
    if (location.pathname.startsWith("/macro/")) return "Редактирование макроса";
    return location.pathname;
  };

  const isMacroEdit =
    location.pathname.startsWith("/macro/") && location.pathname !== "/macro/create";

  const handleDelete = async () => {
    const macroId = location.pathname.split("/").pop();
    if (!macroId) return;

    await deleteMacro(macroId);
    navigate("/");
  };

  return (
    <div className="window">
      <div className="window-title">
        <div className="window-title-text andika" data-draggable>
          <img src={logo} alt="logo" className="w-4 h-4" />
          <span className="cursor-pointer" onClick={() => navigate("/")}>
            repeat what shit
          </span>
          <span className="arrow">
            <IoArrowForward />
          </span>
          <span className="text-xs">{getPageTitle()}</span>
        </div>

        {location.pathname == "/" && (
          <button className="action-btn" onClick={() => navigate("/macro/create")}>
            Cоздать макрос
          </button>
        )}

        {isMacroEdit && (
          <button className="action-btn" onClick={handleDelete}>
            <IoTrash /> Удалить макрос
          </button>
        )}

        <div className="window-title-buttons">
          <button className="close" onClick={() => WindowHide()}>
            <IoClose />
          </button>
        </div>
      </div>

      <div className="window-content">{children}</div>
    </div>
  );
}
