import type {ComponentProps} from 'react';
import {Switch as NativeSwitch} from 'react-native';

export type SwitchProps = ComponentProps<typeof NativeSwitch> & {};
const Switch = ({children, ...props}: SwitchProps) => {
  return <NativeSwitch {...props}>{children}</NativeSwitch>;
};

export default Switch;
