import { DashboardName } from "./meta";
import { AppName } from "@/components/App/meta";
import {
  useAppDispatch,
  withResetAction,
  withAppHeaderAction,
} from "@/components/App/Context";
import {
  useAppI18nDispatch,
  withAppI18nAction,
  withAppI18nResetAction,
} from "@/components/App/I18Next";
import { useEffect } from "react";

const DashboardMain = () => {
  const i18nDispatch = useAppI18nDispatch();
  useEffect(() => {
    i18nDispatch?.(withAppI18nAction({ namespace: DashboardName }));
    return () => i18nDispatch?.(withAppI18nResetAction());
  }, [i18nDispatch]);

  const dispatch = useAppDispatch();
  useEffect(() => {
    dispatch?.(
      withAppHeaderAction({
        title: AppName.toUpperCase(),
      }),
    );
    return () => dispatch?.(withResetAction());
  }, [dispatch]);

  return (
    <>
      <div>Hello, this is my first Page.It name is Home Page 1</div>
      <div>Hello, this is my first Page.It name is Home Page 2</div>
      <div>Hello, this is my first Page.It name is Home Page 3</div>
      <div>Hello, this is my first Page.It name is Home Page 4</div>
      <div>Hello, this is my first Page.It name is Home Page 5</div>
      <div>Hello, this is my first Page.It name is Home Page 6</div>
      <div>Hello, this is my first Page.It name is Home Page 7</div>
      <div>Hello, this is my first Page.It name is Home Page 8</div>
      <div>Hello, this is my first Page.It name is Home Page 9</div>
      <div>Hello, this is my first Page.It name is Home Page 1</div>
      <div>Hello, this is my first Page.It name is Home Page 2</div>
      <div>Hello, this is my first Page.It name is Home Page 3</div>
      <div>Hello, this is my first Page.It name is Home Page 4</div>
      <div>Hello, this is my first Page.It name is Home Page 5</div>
      <div>Hello, this is my first Page.It name is Home Page 6</div>
      <div>Hello, this is my first Page.It name is Home Page 7</div>
      <div>Hello, this is my first Page.It name is Home Page 8</div>
      <div>Hello, this is my first Page.It name is Home Page 9</div>
      <div>Hello, this is my first Page.It name is Home Page 1</div>
      <div>Hello, this is my first Page.It name is Home Page 2</div>
      <div>Hello, this is my first Page.It name is Home Page 3</div>
      <div>Hello, this is my first Page.It name is Home Page 4</div>
      <div>Hello, this is my first Page.It name is Home Page 5</div>
      <div>Hello, this is my first Page.It name is Home Page 6</div>
      <div>Hello, this is my first Page.It name is Home Page 7</div>
      <div>Hello, this is my first Page.It name is Home Page 8</div>
      <div>Hello, this is my first Page.It name is Home Page 9</div>
      <div>Hello, this is my first Page.It name is Home Page 1</div>
      <div>Hello, this is my first Page.It name is Home Page 2</div>
      <div>Hello, this is my first Page.It name is Home Page 3</div>
      <div>Hello, this is my first Page.It name is Home Page 4</div>
      <div>Hello, this is my first Page.It name is Home Page 5</div>
      <div>Hello, this is my first Page.It name is Home Page 6</div>
      <div>Hello, this is my first Page.It name is Home Page 7</div>
      <div>Hello, this is my first Page.It name is Home Page 8</div>
      <div>Hello, this is my first Page.It name is Home Page 9</div>
      <div>Hello, this is my first Page.It name is Home Page 1</div>
      <div>Hello, this is my first Page.It name is Home Page 2</div>
      <div>Hello, this is my first Page.It name is Home Page 3</div>
      <div>Hello, this is my first Page.It name is Home Page 4</div>
      <div>Hello, this is my first Page.It name is Home Page 5</div>
      <div>Hello, this is my first Page.It name is Home Page 6</div>
      <div>Hello, this is my first Page.It name is Home Page 7</div>
      <div>Hello, this is my first Page.It name is Home Page 8</div>
      <div>Hello, this is my first Page.It name is Home Page 9</div>
    </>
  );
};

export default DashboardMain;
