import { AppName } from "@/components/App/meta";
import { useRef, useEffect, type ComponentPropsWithRef } from "react";
import { toDataURL, type QRCodeToDataURLOptions } from "qrcode";

type QRCodeViewProps = {
  value?: string;
  onDataURLChange?: (url: string, value?: string) => void;
  qrCodeOpts?: QRCodeToDataURLOptions;
  className?: string;
} & Pick<ComponentPropsWithRef<"canvas">, "ref" | "className">;
export const QRCodeView = ({
  value,
  onDataURLChange,
  ref,
  qrCodeOpts,
  className,
}: QRCodeViewProps) => {
  const handleDataURLChange = (url: string) => {
    onDataURLChange?.(url, value);
  };

  const internalRef = useRef<HTMLCanvasElement | null>(null);
  const canvasRef = ref && typeof ref !== "function" ? ref : internalRef;

  useEffect(() => {
    if (ref && typeof ref === "function") ref(canvasRef.current);
    if (!canvasRef.current) return;

    const canvas = canvasRef.current;
    toDataURL(canvas, value || "", qrCodeOpts).then(handleDataURLChange);
  }, [value, qrCodeOpts, canvasRef.current]);

  return <canvas className={className} ref={canvasRef} />;
};

type AppQRCodeMeta = {
  name?: string;
  module: string;
  resource: string;
};
type withAppQRCodeParams = ConstructorParameters<typeof URLSearchParams>[0];
export function withAppQRCode<T extends withAppQRCodeParams>(
  params: T,
  meta: AppQRCodeMeta,
) {
  const searchParams = new URLSearchParams(params);
  return `${meta.name || AppName}://${meta.module}/${meta.resource}?${searchParams.toString()}`;
}

export function toAppQRCode(
  value: string,
  meta: AppQRCodeMeta,
): Record<string, string> | undefined {
  if (!URL.canParse(value)) return;

  const url = new URL(value);
  if (url.protocol.slice(0, -1) !== (meta.name || AppName)) return;
  if (url.host !== meta.module) return;
  if (url.pathname.slice(1) !== meta.resource) return;
  return Object.fromEntries(url.searchParams);
}
