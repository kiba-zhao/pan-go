import { ClusterName, PassportName } from "./meta";
import { FileScan } from "./Icon";
import {
  type ComponentProps,
  useState,
  useEffect,
  useRef,
  useMemo,
  type ComponentPropsWithoutRef,
} from "react";
import {
  QRCodeView,
  QRCodeScan,
  QRFileInputAction,
  withAppQRCode,
  toAppQRCode,
} from "@/components/App/QRCode";
import { buttonVariants } from "@/components/ui/button";

type PassportDataType = {
  token: string;
  deviceName: string;
  peerId: string;
};
type PassportQRCodeProps = Omit<
  ComponentProps<typeof QRCodeView>,
  "value" | "ref"
> &
  ComponentPropsWithoutRef<"div">;
export const PassportQRCode = ({
  qrCodeOpts,
  onDataURLChange,
  ...props
}: PassportQRCodeProps) => {
  const value = withAppQRCode(
    {
      token: "123456",
      deviceName: "desktop",
      peerId: "abcdefg123",
    },
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

  const divRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!divRef.current || !observer) return;
    observer.observe(divRef.current);
    return () => observer.disconnect();
  }, [divRef.current, observer]);

  return (
    <div {...props} ref={divRef}>
      <QRCodeView
        className="size-full"
        qrCodeOpts={qrCodeOpts_}
        value={value}
        onDataURLChange={onDataURLChange}
      />
    </div>
  );
};

type PassportQRScanProps = {
  onScan?: (value: PassportDataType) => void;
};
export const PassportQRScan = ({ onScan }: PassportQRScanProps) => {
  const videoRef = useRef<HTMLVideoElement>(null);
  const handleScan = (value: string) => {
    if (!videoRef.current) return;
    const { token, deviceName, peerId } =
      toAppQRCode(value, {
        module: ClusterName,
        resource: PassportName,
      }) || ({} as PassportDataType);
    if (!token || !deviceName || !peerId) return;
    onScan?.({ token, deviceName, peerId });
    videoRef.current.pause();
  };

  const [height, setHeight] = useState(0);
  const observer = useMemo(
    () =>
      new ResizeObserver((entries) => {
        for (const entry of entries) {
          setHeight(entry.contentRect.width);
        }
      }),
    [],
  );
  const canvasRef = useRef<HTMLCanvasElement>(null);
  useEffect(() => {
    if (!canvasRef.current || !observer) return;
    observer.observe(canvasRef.current);
    return () => observer.disconnect();
  }, [canvasRef.current, observer]);

  return (
    <QRCodeScan
      className="w-full"
      width={height}
      height={height}
      onScan={handleScan}
      videoRef={videoRef}
      ref={canvasRef}
    />
  );
};

type PassportQRFileInputActionProps = {
  onChange?: (value: PassportDataType) => void;
};
export const PassportQRFileInputAction = ({
  onChange,
}: PassportQRFileInputActionProps) => {
  const handleChange = (value: string) => {
    const { token, deviceName, peerId } =
      toAppQRCode(value, {
        module: ClusterName,
        resource: PassportName,
      }) || ({} as PassportDataType);
    if (!token || !deviceName || !peerId) return;
    onChange?.({ token, deviceName, peerId });
  };
  return (
    <>
      <label
        htmlFor="passport-file-input"
        className={buttonVariants({ variant: "ghost", size: "sm" })}
      >
        <FileScan />
        读取二维码文件
      </label>
      <QRFileInputAction onChange={handleChange} id="passport-file-input" />
    </>
  );
};
