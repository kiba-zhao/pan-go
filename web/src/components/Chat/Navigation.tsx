import { NavList } from "@/components/App/Navigation";
import type { ComponentProps } from "react";

const ChatNavigation = ({ className }: ComponentProps<typeof NavList>) => {
  // TODO: Chat Navigation
  return <NavList className={className}>Chat Navigation</NavList>;
};

export default ChatNavigation;
