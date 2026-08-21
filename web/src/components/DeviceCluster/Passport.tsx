import { ClusterName, PassportName } from "./meta";
import {
  type ComponentProps,
  useState,
  useEffect,
  useRef,
  useMemo,
  type ComponentPropsWithoutRef,
} from "react";
import { QRCodeView, withAppQRCode } from "@/components/App/QRCode";

type PassportQRCode = Omit<ComponentProps<typeof QRCodeView>, "value" | "ref"> &
  ComponentPropsWithoutRef<"div">;
export const PassportQRCode = ({
  qrCodeOpts,
  onDataURLChange,
  ...props
}: PassportQRCode) => {
  const value = withAppQRCode(
    { token: "123456", timeout: Date.now().toString(), deviceName: "desktop" },
    { module: ClusterName, resource: PassportName },
  );

  const [qrCodeOpts_, setQRCodeOpts] =
    useState<ComponentProps<typeof QRCodeView>["qrCodeOpts"]>(qrCodeOpts);

  const observer = useMemo(
    () =>
      !qrCodeOpts?.width &&
      new ResizeObserver((entries) => {
        for (const entry of entries) {
          setQRCodeOpts({
            ...(qrCodeOpts || {}),
            width: entry.contentRect.width,
          });
        }
      }),
    [qrCodeOpts?.width],
  );

  const devRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!devRef.current || !observer) return;
    observer.observe(devRef.current);
    return () => observer.disconnect();
  }, [devRef.current, observer]);

  return (
    <div {...props} ref={devRef}>
      <QRCodeView
        className="size-full"
        qrCodeOpts={qrCodeOpts_}
        value={value}
        onDataURLChange={onDataURLChange}
      />
    </div>
  );
};
