import debounce from 'lodash.debounce';
import {Fragment, useCallback, useEffect, useMemo, useState} from 'react';
import {GestureResponderEvent, View} from 'react-native';
import Button from './Button';
import {useTranslation} from './I18Next';
import type {IconProps} from './Icon';
import Icon from './Icon';
import {RowLayout} from './Layout';
import Paper from './Paper';
import type {ScreenModalProps} from './ScreenBase';
import {ScreenModal} from './ScreenBase';
import {scale} from './SizeMatters';
import {default as Switch, type SwitchProps} from './Switch';
import Text from './Text';
import type {TextInputProps} from './TextInput';
import TextInput from './TextInput';
import {useTheme} from './Theme';

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

export type TextFieldModalProps = {
  onValid?: (text?: string, value?: string) => boolean | Error;
  onSubmit?: (text: string) => void;
  labelText?: string;
  labelI18nKey?: string;
  okText?: string;
  okI18nKey?: string;
  cancelText?: string;
  cancelI18nKey?: string;
  placeholder?: string;
  placeholderI18nKey?: string;
  helperText?: string;
  helperI18nKey?: string;
  validMode?: 'strict' | 'inclusive';
} & Pick<ScreenModalProps, 'onClose' | 'visible'> &
  Omit<TextInputProps, 'onChangeText' | 'onSubmitEditing'>;
export const TextFieldModal = ({
  visible,
  value,
  onValid,
  onClose,
  onSubmit,
  placeholder,
  placeholderI18nKey,
  helperText,
  helperI18nKey,
  labelText,
  labelI18nKey,
  okText,
  okI18nKey = 'action.ok',
  cancelText,
  cancelI18nKey = 'action.cancel',
  validMode = 'strict',
  ...props
}: TextFieldModalProps) => {
  const {t} = useTranslation();
  const {colors, sizes} = useTheme();
  const [text, setText] = useState(value);
  const [validState, setValidState] = useState<boolean | Error>(false);

  useEffect(() => {
    setText(value);
    setValidState(false);
  }, [value]);

  const handleValid = useCallback(
    (text: string, value?: string) => {
      let state = onValid?.(text, value);
      if (state === void 0) {
        state = text.length > 0 && text !== value;
      }
      return state;
    },
    [onValid],
  );

  const setValidStateWithDelay = useCallback(debounce(setValidState, 500), [
    setValidState,
  ]);

  const handleValidStateChange = useCallback(
    debounce((text: string, value?: string) => {
      const validState = handleValid(text, value);
      setValidState(validState);
      return validState;
    }, 500),
    [setValidState, handleValid],
  );

  const handleChangeText = (text: string) => {
    if (validMode === 'inclusive') {
      handleValidStateChange(text, value);
    } else {
      const validState = handleValid(text, value);
      setValidStateWithDelay(validState);
      if (validState instanceof Error) return;
    }
    setText(text);
  };

  const handleSubmit = (event?: GestureResponderEvent) => {
    const validState = handleValid(text || '', value);
    setValidState(validState);
    if (validState !== true) return;
    onSubmit?.(text || '');
  };

  const handleCancel = (event: GestureResponderEvent) => {
    onClose?.(event);
    setText(value);
    setValidState(false);
  };

  const title = useMemo(() => {
    if (labelText) return labelText;
    if (labelI18nKey) return t(labelI18nKey);
    return void 0;
  }, [labelText, labelI18nKey]);

  const placeholder_ = useMemo(() => {
    if (placeholder) return placeholder;
    if (placeholderI18nKey) return t(placeholderI18nKey);
    return void 0;
  }, [placeholder, placeholderI18nKey]);

  const helperText_ = useMemo(() => {
    if (helperText) return helperText;
    if (helperI18nKey) return t(helperI18nKey);
    return void 0;
  }, [helperText, helperI18nKey]);

  return (
    <ScreenModal visible={visible} onClose={handleCancel}>
      <Paper padding={1.5} style={{width: '100%'}}>
        {title && (
          <Text size="title" font="medium">
            {title}
          </Text>
        )}
        <TextInput
          style={{
            width: '100%',
            borderColor: colors.divider,
            borderBottomWidth: 1,
          }}
          value={text}
          onChangeText={handleChangeText}
          onSubmitEditing={event => handleSubmit()}
          placeholder={placeholder_}
          {...props}
        />
        <Text padding={[0, 0.5]} color="textSecondary" size="small">
          {helperText_}
        </Text>
        <Paper bgColor="transparent" padding={[3, 0, 1, 0]}>
          <RowLayout
            style={{justifyContent: 'flex-end', gap: scale(sizes.base * 2)}}>
            <View style={{flex: 1}}>
              <Button
                color={colors.secondary}
                onPress={handleSubmit}
                title={okText || t(okI18nKey)}
                disabled={validState !== true}
              />
            </View>
            <View style={{flex: 1}}>
              <Button
                onPress={handleCancel}
                title={cancelText || t(cancelI18nKey)}
              />
            </View>
          </RowLayout>
        </Paper>
      </Paper>
    </ScreenModal>
  );
};
