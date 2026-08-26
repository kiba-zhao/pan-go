import { AppName } from "./meta";
import { useTranslation, I18nVariant, useAppI18n } from "./I18Next";
import {
  useRef,
  useEffect,
  useState,
  type ComponentProps,
  type ComponentPropsWithRef,
  type ChangeEvent,
  SyntheticEvent,
} from "react";
import { default as jsQR, QRCode as JSQRCode } from "jsqr";
import { toDataURL, type QRCodeToDataURLOptions } from "qrcode";
import { useBrowser } from "@/components/App/Browser";
import { buttonVariants } from "@/components/ui/button";
import { Alert, AlertTitle, AlertDescription } from "@/components/ui/alert";
import { cn } from "@/lib/utils";

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

type QRCodeScanProps = {
  onScan?: (value: string) => void;
  videoRef?: CameraVideoProps["ref"];
} & ComponentPropsWithRef<"canvas">;
export const QRCodeScan = ({
  ref,
  videoRef,
  className,
  onScan,
  ...props
}: QRCodeScanProps) => {
  const { namespace } = useAppI18n();
  const { t } = useTranslation(namespace);
  const internalRef = useRef<HTMLCanvasElement | null>(null);
  const canvasRef = ref && typeof ref !== "function" ? ref : internalRef;
  const frameRef = useRef<number | null>(null);
  const browser = useBrowser();

  const handleVideoPlay = (video: HTMLVideoElement) => {
    const window = browser?.window;
    if (!window) {
      video.pause();
      return;
    }
    if (video.readyState === video.HAVE_ENOUGH_DATA && canvasRef.current) {
      const canvas = canvasRef.current;
      const ctx = canvas.getContext("2d");
      if (ctx) {
        ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
        const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
        const results = jsQR(imageData.data, imageData.width, imageData.height);
        if (results) {
          drawDetectionBox(ctx, results.location);
          onScan?.(results.data || "");
        }
      }
    }
    if (video.paused) return;
    frameRef.current = window.requestAnimationFrame(() =>
      handleVideoPlay(video),
    );
  };

  const handleVideoPause = () => {
    if (frameRef.current) {
      window.cancelAnimationFrame(frameRef.current);
      frameRef.current = null;
    }
  };

  const [error, setError] = useState<Error | null>(null);
  const handleVideoError = (event: SyntheticEvent) => {
    const { error } = event.nativeEvent as ErrorEvent;
    setError(error);
  };

  return (
    <div className="relative">
      <Alert
        variant="destructive"
        className={cn(
          "absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 bg-transparent border-none",
          !error && "hidden",
        )}
      >
        <AlertTitle className="text-center">
          {t(`${I18nVariant.Error}.QRCodeScanError`)}
        </AlertTitle>
        <AlertDescription className="text-center">
          {t(`${I18nVariant.Error}.QRCodeScanErrors.${error?.name}`, {
            defaultValue: error?.message,
          })}
        </AlertDescription>
      </Alert>
      <canvas
        {...props}
        className={cn(className, error && "bg-muted")}
        ref={canvasRef}
      />
      <CameraVideo
        ref={videoRef}
        onPlay={(event) => handleVideoPlay(event.target as HTMLVideoElement)}
        onPause={handleVideoPause}
        onError={handleVideoError}
        muted
        playsInline
        hidden
      />
    </div>
  );
};

type QRCodeLocation = JSQRCode["location"];
type QRCodePoint = QRCodeLocation["topLeftCorner"];

type CameraVideoProps = ComponentPropsWithRef<"video">;
const CameraVideo = ({ ref, ...props }: CameraVideoProps) => {
  const internalRef = useRef<HTMLVideoElement | null>(null);
  const videoRef = ref && typeof ref !== "function" ? ref : internalRef;
  const browser = useBrowser();

  useEffect(() => {
    if (!browser || !videoRef.current) return;
    const video = videoRef.current;
    const window = browser.window;
    const { navigator } = window;

    const promise = navigator.mediaDevices
      .getUserMedia({ video: { facingMode: "environment" } })
      .then((stream) => {
        video.srcObject = stream;
        video.play();
        return stream;
      });

    promise.catch((error) => {
      video.dispatchEvent(new ErrorEvent("error", { error }));
    });

    return () => {
      promise.then(destructorForCameraVideo);
    };
  }, [videoRef.current, browser]);

  return <video {...props} ref={videoRef}></video>;
};

type QRFileInputActionProps = {
  onChange: (value: string) => void;
} & Pick<ComponentProps<"input">, "id" | "accept">;
export const QRFileInputAction = ({
  onChange,
  id,
  accept = "image/*",
}: QRFileInputActionProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const imageRef = useRef<HTMLImageElement | null>(null);
  const [image, setImage] = useState<string | undefined>();

  const onFileChange = (e: ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setImage(URL.createObjectURL(file));
    e.target.value = "";
  };

  const onLoad = () => {
    if (!canvasRef.current || !imageRef.current) return;
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d", { willReadFrequently: true });
    const img = imageRef.current;
    canvas.width = img.width;
    canvas.height = img.height;
    if (!ctx) return;

    ctx.reset();
    ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const results = jsQR(imageData.data, imageData.width, imageData.height);
    onChange(results?.data || "");
  };

  return (
    <>
      <input
        type="file"
        hidden
        accept={accept}
        id={id}
        onChange={onFileChange}
      />
      <img
        hidden
        src={image}
        alt="qr image decode"
        onLoad={onLoad}
        ref={imageRef}
      />
      <canvas ref={canvasRef} hidden />
    </>
  );
};

type QRCodeSaveActionProps = {
  variant?: NonNullable<Parameters<typeof buttonVariants>[0]>["variant"];
  size?: NonNullable<Parameters<typeof buttonVariants>[0]>["size"];
} & Pick<ComponentProps<"a">, "className" | "children" | "href" | "download">;
export const QRCodeSaveAction = ({
  variant = "secondary",
  size = "sm",
  className,
  children,
  ...props
}: QRCodeSaveActionProps) => {
  return (
    <a {...props} className={cn(buttonVariants({ variant, size }), className)}>
      {children}
    </a>
  );
};

function destructorForCameraVideo(stream: MediaStream) {
  stream.getTracks().forEach((track) => track.stop());
}

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

export function toAppQRCode<T extends Record<string, string>>(
  value: string,
  meta: AppQRCodeMeta,
): T | undefined {
  if (!URL.canParse(value)) return;

  const url = new URL(value);
  if (url.protocol.slice(0, -1) !== (meta.name || AppName)) return;
  if (url.host !== meta.module) return;
  if (url.pathname.slice(1) !== meta.resource) return;
  return Object.fromEntries(url.searchParams) as T;
}

function drawDetectionBox(
  canvas: CanvasRenderingContext2D,
  location: QRCodeLocation,
  color?: CanvasFillStrokeStyles["strokeStyle"],
) {
  const color_ = color || "red";
  drawLine(canvas, location.topLeftCorner, location.topRightCorner, color_);
  drawLine(canvas, location.topRightCorner, location.bottomRightCorner, color_);
  drawLine(
    canvas,
    location.bottomRightCorner,
    location.bottomLeftCorner,
    color_,
  );
  drawLine(canvas, location.bottomLeftCorner, location.topLeftCorner, color_);
}

function drawLine(
  canvas: CanvasRenderingContext2D,
  begin: QRCodePoint,
  end: QRCodePoint,
  color: CanvasFillStrokeStyles["strokeStyle"],
) {
  canvas.beginPath();
  canvas.moveTo(begin.x, begin.y);
  canvas.lineTo(end.x, end.y);
  canvas.lineWidth = 4;
  canvas.strokeStyle = color;
  canvas.stroke();
}
