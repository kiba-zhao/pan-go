import { type ElementType } from "react";

export type AsProps<Props> = Props extends {
  as?: any;
}
  ? never
  : { as: ElementType<Props> } & Props;
