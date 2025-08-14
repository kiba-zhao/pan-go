import type {ModalProps, ViewStyle} from 'react-native';
import {
  ActivityIndicator,
  Modal as NativeModal,
  RefreshControl,
  SafeAreaView,
  ScrollView,
  StyleSheet,
  TouchableWithoutFeedback,
  View,
} from 'react-native';
import {scale} from './SizeMatters';
import {withTheme as withStyleTheme} from './StyleSheet';

import {useCallback, type ComponentProps} from 'react';
import type {NativeSyntheticEvent} from 'react-native';
import type {PressableProps, PressableStateCallbackType} from './Pressable';
import Pressable from './Pressable';
import Text from './Text';
import {useTheme} from './Theme';

export type ScreenLayoutProps = ComponentProps<typeof ScrollView>;
export const ScreenLayout = ({
  children,
  style,
  ...props
}: ScreenLayoutProps) => {
  const {sizes} = useTheme();

  return (
    <SafeAreaView>
      <ScrollView
        {...props}
        style={StyleSheet.compose(
          {paddingHorizontal: scale(sizes.base)},
          style,
        )}>
        {children}
      </ScrollView>
    </SafeAreaView>
  );
};

export const ScreenLoading = () => {
  return (
    <View
      style={{
        justifyContent: 'center',
        alignItems: 'center',
        height: '100%',
        width: '100%',
        zIndex: 1,
        position: 'absolute',
      }}>
      <ActivityIndicator size={60} />
    </View>
  );
};

type NativeRefreshControlProps = ComponentProps<typeof RefreshControl>;
export type ScreenRefreshControlProps = Partial<
  Pick<NativeRefreshControlProps, 'refreshing'>
> &
  Omit<NativeRefreshControlProps, 'refreshing'>;
export const ScreenRefreshControl = ({
  refreshing,
  ...props
}: ScreenRefreshControlProps) => {
  return <RefreshControl refreshing={refreshing || false} {...props} />;
};

export const ScreenSafetyFooter = withStyleTheme(
  View,
  ({sizes}) =>
    ({
      height: scale(sizes.base * 2),
    } as ViewStyle),
);

export type ScreenModalProps = {
  onClose?: PressableProps['onPress'];
  containerProps?: Omit<PressableProps, 'children'>;
} & Omit<ModalProps, 'children'> &
  Pick<PressableProps, 'children'>;
export const ScreenModal = ({
  containerProps = {},
  onClose,
  children,
  backdropColor = 'transparent',
  onRequestClose,
  animationType = 'slide',
  ...props
}: ScreenModalProps) => {
  const handleRequestClose = useCallback(
    (event: NativeSyntheticEvent<any>) => {
      (onRequestClose || onClose)?.(event);
    },
    [onRequestClose, onClose],
  );

  const {style, onPress, ...containerProps_} = containerProps;

  const handlePress = useCallback(
    (event: NativeSyntheticEvent<any>) => {
      console.log(`handlePress in ScreenModal`);
      (onPress || onClose)?.(event);
    },
    [onPress, onClose],
  );

  const style_ = useCallback(
    (state: PressableStateCallbackType) => {
      const customStyle = typeof style === 'function' ? style(state) : style;
      return StyleSheet.compose(
        {
          justifyContent: 'flex-end',
          alignItems: 'center',
          height: '100%',
          width: '100%',
        },
        customStyle,
      );
    },
    [style],
  );

  const children_ = useCallback(
    (state: PressableStateCallbackType) => {
      return typeof children === 'function' ? (
        <TouchableWithoutFeedback>{children(state)}</TouchableWithoutFeedback>
      ) : (
        <TouchableWithoutFeedback>{children}</TouchableWithoutFeedback>
      );
    },
    [children],
  );

  return (
    <NativeModal
      onRequestClose={handleRequestClose}
      backdropColor={backdropColor}
      animationType={animationType}
      visible={false}
      {...props}>
      <Pressable style={style_} onPress={handlePress} {...containerProps_}>
        {children_}
      </Pressable>
    </NativeModal>
  );
};

type ScreenEmptyProps = ComponentProps<typeof View>;
export const ScreenEmpty = ({children, style, ...props}: ScreenEmptyProps) => {
  const {sizes} = useTheme();

  const style_ = style || {
    alignItems: 'center',
    backgroundColor: 'transparent',
    paddingVertical: scale(sizes.base * 8),
  };

  return (
    <View style={style_} {...props}>
      {typeof children === 'string' ? (
        <Text font="bold" color="textDisabled" size="title">
          {children}
        </Text>
      ) : (
        children
      )}
    </View>
  );
};
