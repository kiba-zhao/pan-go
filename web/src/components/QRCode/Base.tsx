import type { ChangeEvent, ReactNode } from "react";
import { Fragment, useCallback, useEffect, useRef, useState } from "react";

import type { QRCode as JSQRCode } from "jsqr";
import jsQR from "jsqr";
import type { QRCodeRenderersOptions } from "qrcode";
import { toCanvas } from "qrcode";
import { useBrower } from "../Global/Brower";
import {
  DEFAULT_BASE,
  TranslateProvider,
  useTranslate,
} from "../Global/Translation";

import Button from "@mui/material/Button";
import type { Dispatch } from "react";
import { createContext, useContext, useReducer } from "react";

import PhotoIcon from "@mui/icons-material/Photo";

type QRState = {
  download?: string;
  value?: string;
  url?: string;
};
type QRActionType = "SET";
type QRAction = { type: QRActionType } & QRState;
const QRReducer = (state: QRState, action: QRAction) => {
  const { type, ...state_ } = action;
  switch (type) {
    case "SET":
      return { ...state, ...state_ };
    default:
      return state;
  }
};

const QRContext = createContext<QRState>(null!);
const QRDispatchContext = createContext<Dispatch<QRAction>>(null!);

export const useQR = () => useContext(QRContext);
export const useQRDispatch = () => useContext(QRDispatchContext);
export const QRProvider = ({
  children,
  base = DEFAULT_BASE,
}: {
  children: ReactNode;
  base?: string;
}) => {
  const [state, dispatch] = useReducer(QRReducer, {} as QRState);
  return (
    <TranslateProvider base={base}>
      <QRContext.Provider value={state}>
        <QRDispatchContext.Provider value={dispatch}>
          {children}
        </QRDispatchContext.Provider>
      </QRContext.Provider>
    </TranslateProvider>
  );
};

const QRCodeErrorColor = {
  light: "#ff0000",
};
export type QRCodeProps = {
  value?: string;
  name?: string;
} & QRCodeRenderersOptions;
export const QRCode = ({ name, value, ...opts }: QRCodeProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const dispatch = useQRDispatch();

  useEffect(() => {
    const canvas = canvasRef.current;

    toCanvas(canvas, value || "invalid", {
      ...opts,
      color: !value ? QRCodeErrorColor : opts.color,
    }).then(() => {
      const url = canvas?.toDataURL();
      dispatch({
        type: "SET",
        download: `${name}.png`,
        value: value,
        url,
      });
    });
  }, [name, value, opts]);

  return <canvas ref={canvasRef}></canvas>;
};

export const QRCodeDownloadButton = () => {
  const t = useTranslate();
  const { download, value, url } = useQR();
  return (
    <Button
      variant="contained"
      size="small"
      href={url || ""}
      download={download}
      disabled={!value || value.length <= 0}
      startIcon={<PhotoIcon />}
    >
      {t("button.save")}
    </Button>
  );
};

type CameraVideoProps = {
  onChanged: (video: HTMLVideoElement) => void;
  hidden?: boolean;
};
const CameraVideo = ({ onChanged, hidden }: CameraVideoProps) => {
  const videoRef = useRef<HTMLVideoElement | null>(null);
  const brower = useBrower();
  const cancelRef = useRef<number | null>(null);

  const destructor = useCallback((stream: MediaStream) => {
    stream.getTracks().forEach((track) => track.stop());
    if (!brower || !videoRef.current) return;
    const window = brower.window;
    if (cancelRef.current) {
      window.cancelAnimationFrame(cancelRef.current);
      cancelRef.current = null;
    }
  }, []);

  const next = useCallback(() => {
    if (!brower || !videoRef.current) return;
    const video = videoRef.current;
    const window = brower.window;

    const cancelId = window.requestAnimationFrame(() => {
      if (video.readyState === video.HAVE_ENOUGH_DATA) {
        onChanged(video);
        if (video.paused) {
          return;
        }
      }
      next();
    });
    cancelRef.current = cancelId;
  }, []);

  useEffect(() => {
    if (!brower || !videoRef.current) return;
    const video = videoRef.current;
    const window = brower.window;
    const { navigator } = window;

    const promise = navigator.mediaDevices
      .getUserMedia({ video: { facingMode: "environment" } })
      .then((stream) => {
        video.srcObject = stream;
        video.setAttribute("playsinline", true.toString()); // required to tell iOS safari we don't want fullscreen
        video.play();
        next();
        return stream;
      });

    return () => {
      promise.then(destructor);
    };
  }, []);

  return <video ref={videoRef} hidden={hidden}></video>;
};

export type QRScanChangedEvent = {
  value: string;
  invalid?: boolean;
};

type QRScanProps = {
  onChanged: (event: QRScanChangedEvent) => void;
  once?: boolean;
  width?: string | number;
  height?: string | number;
};
export const QRScan = ({ onChanged, once, width, height }: QRScanProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const valueRef = useRef<string | null>(null);
  const onCameraVideoChanged = (video: HTMLVideoElement) => {
    if (!canvasRef.current) return;
    if (video.paused) return;
    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d", { willReadFrequently: true });
    if (!ctx) return;
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const results = jsQR(imageData.data, imageData.width, imageData.height);
    if (!results) return true;
    drawDetectionBox(ctx, results.location);
    if (valueRef.current === results.data) return;
    const event: QRScanChangedEvent = { value: results.data };
    onChanged(event);
    if (once && event.invalid === false) {
      video.pause();
    }
    return;
  };

  return (
    <>
      <canvas ref={canvasRef} width={width} height={height} />
      <CameraVideo onChanged={onCameraVideoChanged} hidden />
    </>
  );
};

type QRFileScanProps = {
  onChange: (value: string) => void;
};
export const QRFileScan = ({ onChange }: QRFileScanProps) => {
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
    <Fragment>
      <input type="file" hidden accept="image/*" onChange={onFileChange} />
      <img
        hidden
        src={image}
        alt="qr image decode"
        onLoad={onLoad}
        ref={imageRef}
      />
      <canvas ref={canvasRef} hidden />
    </Fragment>
  );
};

type QRCodeLocation = JSQRCode["location"];
type QRCodePoint = QRCodeLocation["topLeftCorner"];

function drawDetectionBox(
  canvas: CanvasRenderingContext2D,
  location: QRCodeLocation,
  color?: CanvasFillStrokeStyles["strokeStyle"]
) {
  const color_ = color || "red";
  drawLine(canvas, location.topLeftCorner, location.topRightCorner, color_);
  drawLine(canvas, location.topRightCorner, location.bottomRightCorner, color_);
  drawLine(
    canvas,
    location.bottomRightCorner,
    location.bottomLeftCorner,
    color_
  );
  drawLine(canvas, location.bottomLeftCorner, location.topLeftCorner, color_);
}

function drawLine(
  canvas: CanvasRenderingContext2D,
  begin: QRCodePoint,
  end: QRCodePoint,
  color: CanvasFillStrokeStyles["strokeStyle"]
) {
  canvas.beginPath();
  canvas.moveTo(begin.x, begin.y);
  canvas.lineTo(end.x, end.y);
  canvas.lineWidth = 4;
  canvas.strokeStyle = color;
  canvas.stroke();
}
