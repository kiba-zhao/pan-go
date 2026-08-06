import { type ComponentProps } from "react";
export const Blank = ({ children, ...props }: ComponentProps<"p">) => (
  <p {...props}>{children || "\u00A0"}</p>
);
