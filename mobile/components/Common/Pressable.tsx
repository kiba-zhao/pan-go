import type {PressableProps, ViewStyle} from 'react-native';
import {Pressable, StyleSheet} from 'react-native';
import {withStyleFunction} from './StyleSheet';

type PressableStyle = PressableProps['style'];
type PressableStyleFunction = Extract<PressableStyle, Function>;
type PressableStyleParamters = Parameters<PressableStyleFunction>;
type testType = PressableStyleParamters[0];

const styles = StyleSheet.create({
  opacity: {
    opacity: 0.6,
  },
});

export const OpacityPressable = withStyleFunction<
  ViewStyle,
  PressableStyleParamters,
  PressableProps
>(Pressable, ({pressed}) => (pressed ? styles.opacity : {}));
