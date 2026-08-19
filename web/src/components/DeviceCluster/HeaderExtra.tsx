import { GoBack } from "@/components/App/HeaderExtra";
import { useLocation } from "@/components/App/Router";

export type ClusterGoBackState = {
  to: string;
  text?: string;
};
export const ClusterHeaderExtra = () => {
  const location = useLocation();
  const { to, text } = (location.state?.goback || {}) as ClusterGoBackState;

  if (!to) {
    return null;
  }
  return <GoBack to={to}>{text}</GoBack>;
};

export function withGoBackState(state: ClusterGoBackState) {
  return { goback: state };
}
