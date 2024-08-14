import { useCallback, useEffect, useRef } from "react";

import jsQR from "jsqr";
import type { QRCodeRenderersOptions } from "qrcode";
import { toCanvas } from "qrcode";
import { useBrower } from "./Brower";

const QRCodeErrorColor = {
  light: "#ff0000",
};
export type QRCodeProps = { value: string | null } & QRCodeRenderersOptions;
export const QRCode = ({ value, ...opts }: QRCodeProps) => {
  const canvasRef = useRef<HTMLCanvasElement | null>(null);

  useEffect(() => {
    toCanvas(canvasRef.current, value || "null", {
      ...opts,
      color: value === null ? QRCodeErrorColor : opts.color,
    });
  }, [value, opts]);

  return <canvas ref={canvasRef}></canvas>;
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
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.drawImage(video, 0, 0, canvas.width, canvas.height);
    const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
    const results = jsQR(imageData.data, imageData.width, imageData.height);
    if (!results) return true;
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
