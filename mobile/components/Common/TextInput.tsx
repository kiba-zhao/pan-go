import type {ComponentProps} from 'react';
import {TextInput as NativeTextInput} from 'react-native';

type NativeTextInputProps = ComponentProps<typeof NativeTextInput>;
export type TextInputProps = {} & NativeTextInputProps;

const TextInput = ({...props}: TextInputProps) => {
  return <NativeTextInput {...props} />;
};

export default TextInput;
