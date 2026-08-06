import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";
import { Address4, Address6 } from "ip-address";

export function cn(...classNames: ClassValue[]) {
  return twMerge(clsx(classNames));
}

export type { ClassValue as ClassName };

export function newIPAddr(addr: string): undefined | Address4 | Address6 {
  if (Address4.isValid(addr)) {
    return new Address4(addr);
  }
  if (Address6.isValid(addr)) {
    return new Address6(addr);
  }
}
