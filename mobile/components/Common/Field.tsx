import {Fragment, useEffect, useState} from 'react';
import Button from './Button';
import type {IconProps} from './Icon';
import Icon from './Icon';
import {RowLayout} from './Layout';
import Paper from './Paper';
import type {ScreenModalProps} from './ScreenBase';
import {ScreenModal} from './ScreenBase';
import {default as Switch, type SwitchProps} from './Switch';
import Text from './Text';
import type {TextInputProps} from './TextInput';
import TextInput from './TextInput';

type TextFieldItemProps = {
  label: string;
  text?: string;
  editable?: boolean;
};
export const TextFieldItem = ({label, text, editable}: TextFieldItemProps) => {
  return (
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={text == void 0 ? 'textDisabled' : 'textPrimary'}>
        {label}
      </Text>
      <RowLayout style={{alignItems: 'center'}}>
        <Text
          size="small"
          color="textSecondary"
          margin={editable ? [0, 0.5, 0, 0] : 0}>
          {text || ''}
        </Text>
        {editable && <Icon name="create-outline" color="textPrimary" />}
      </RowLayout>
    </Fragment>
  );
};

type SwitchFieldItemProps = {
  label: string;
  value?: boolean;
} & Omit<SwitchProps, 'value' | 'disabled'>;
export const SwitchFieldItem = ({
  label,
  value,
  ...props
}: SwitchFieldItemProps) => {
  return (
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={value === void 0 ? 'textDisabled' : 'textPrimary'}>
        {label}
      </Text>
      <Switch {...props} value={!!value} disabled={value === void 0} />
    </Fragment>
  );
};

type IconFieldItemProps = {
  label: string;
} & Omit<IconProps, 'label'>;
export const IconFieldItem = ({
  label,
  disabled,
  ...props
}: IconFieldItemProps) => {
  return (
    <Fragment>
      <Text
        size="small"
        style={{flex: 1}}
        color={disabled ? 'textDisabled' : 'textPrimary'}>
        {label}
      </Text>
      <Icon
        color={disabled ? 'textDisabled' : 'textPrimary'}
        {...props}
        disabled={disabled}
      />
    </Fragment>
  );
};

type TextFieldModalProps = {
  title: string;
  onSubmit?: (text?: string) => void;
  submitLabel?: string;
} & Pick<ScreenModalProps, 'onClose' | 'visible'> &
  Omit<TextInputProps, 'onChangeText' | 'onSubmitEditing'>;
export const TextFieldModal = ({
  title,
  visible,
  value,
  onClose,
  onSubmit,
  submitLabel,
  ...props
}: TextFieldModalProps) => {
  const [text, setText] = useState(value);

  useEffect(() => {
    setText(value);
  }, [value]);

  const handleSubmit = () => {
    onSubmit?.(text);
  };
  return (
    <ScreenModal visible={visible} onClose={onClose}>
      <RowLayout>
        <Text style={{flex: 1}}>{title}</Text>
        <Icon name="close-outline" onPress={onClose || void 0} />
      </RowLayout>
      <TextInput
        style={{width: '100%', borderColor: 'black', borderBottomWidth: 1}}
        value={text}
        onChangeText={setText}
        onSubmitEditing={handleSubmit}
        {...props}
      />
      <Paper bgColor="transparent" padding={[3, 0, 1, 0]}>
        <Button
          onPress={handleSubmit}
          title={submitLabel || 'Submit'}
          disabled={!text}
        />
      </Paper>
    </ScreenModal>
  );
};
