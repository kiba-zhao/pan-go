import { useAppDispatch, resetToBlank } from "../App/Context";
import { useEffect } from "react";

const DashboardMain = () => {
  const dispatch = useAppDispatch();

  useEffect(() => {
    dispatch?.({
      header: { title: import.meta.env.VITE_APP_NAME?.toUpperCase() },
    });
    return () => dispatch?.(resetToBlank());
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
