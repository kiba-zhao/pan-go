import type {ComponentProps} from 'react';
import {Button as NativeButton} from 'react-native';

type ButtonProps = ComponentProps<typeof NativeButton>;
const Button = ({...props}: ButtonProps) => {
  return <NativeButton {...props} />;
};

export default Button;
