import type {ModalProps, ViewStyle} from 'react-native';
import {
  ActivityIndicator,
  Modal as NativeModal,
  Pressable as NativePressable,
  RefreshControl,
  SafeAreaView,
  ScrollView,
  StyleSheet,
  View,
} from 'react-native';
import {scale} from './SizeMatters';
import {withTheme as withStyleTheme} from './StyleSheet';

import type {ComponentProps} from 'react';
import type {GestureResponderEvent, NativeSyntheticEvent} from 'react-native';
import type {PressableProps} from './Pressable';
import Pressable from './Pressable';
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
  onClose,
  containerProps = {},
  children,
  backdropColor = 'transparent',
  onRequestClose,
  ...props
}: ScreenModalProps) => {
  const {sizes} = useTheme();
  const handleRequestClose = (event: NativeSyntheticEvent<any>) => {
    onRequestClose?.(event);
    onClose?.(event);
  };

  const handleBreakPress = (event: GestureResponderEvent) => {
    event.stopPropagation();
    containerProps.onPress?.(event);
  };
  return (
    <NativeModal
      onRequestClose={handleRequestClose}
      backdropColor={backdropColor}
      animationType="fade"
      visible={false}
      {...props}>
      <NativePressable
        style={{
          justifyContent: 'center',
          alignItems: 'center',
          height: '100%',
          width: '100%',
          paddingHorizontal: scale(sizes.base * 2),
        }}
        onPress={onClose}>
        <Pressable
          bgColor={'surface'}
          radius={1}
          padding={1.5}
          style={{width: '100%', maxWidth: 400}}
          {...containerProps}
          onPress={handleBreakPress}>
          {children}
        </Pressable>
      </NativePressable>
    </NativeModal>
  );
};
