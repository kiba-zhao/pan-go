import { Button } from "@/components/ui/button";
import { Separator as ShadcnSeparator } from "@/components/ui/separator";
import { useMedia } from "@/lib/hooks";
import {
  useAppContext,
  useAppDispatch,
  AppAsideMode,
  withAppAction,
} from "./Context";
import { PanelLeft } from "./Icon";

export const AsideModeControl = () => {
  const { asideMode } = useAppContext();
  const dispatch = useAppDispatch();

  const isWide = useMedia("(min-width: 48rem)");

  const handleClick = () => {
    let newAsideMode;

    switch (asideMode) {
      case AppAsideMode.Collapsed:
        newAsideMode = isWide ? AppAsideMode.Hidden : void 0;
        break;
      case AppAsideMode.Hidden:
        newAsideMode = void 0;
        break;
      default:
        newAsideMode = AppAsideMode.Collapsed;
        break;
    }
    dispatch?.(withAppAction({ asideMode: newAsideMode }));
  };
  return (
    <Button size="icon" variant="ghost" onClick={handleClick}>
      <PanelLeft />
    </Button>
  );
};

export const Separator = () => {
  return (
    <ShadcnSeparator
      orientation="vertical"
      className="mr-2 data-[orientation=vertical]:h-4"
    />
  );
};
