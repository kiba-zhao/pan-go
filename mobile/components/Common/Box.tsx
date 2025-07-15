import {useMemo} from 'react';
import type {PressableProps, ViewProps, ViewStyle} from 'react-native';
import {Pressable, StyleSheet, View} from 'react-native';
import type {PaperColor} from './Paper';
import {
  transformBorderColor,
  transformBorderWidthStyle,
  transformGapStyle,
  transformMarginStyle,
  transformPaddingStyle,
  transformRadiusStyle,
} from './StyleSheet';
import {transformColor, useTheme} from './Theme';

export type BoxPressableProps = Pick<
  PressableProps,
  'onPress' | 'onLongPress' | 'onPressIn' | 'onPressOut'
>;

type BoxColor = PaperColor;
export type BoxProps = BoxPressableProps &
  Omit<ViewProps, 'children'> & {
    bgColor?: BoxColor;
    borderColor?:
      | BoxColor
      | [BoxColor, BoxColor]
      | [BoxColor, BoxColor, BoxColor, BoxColor];
    padding?: number | [number, number] | [number, number, number, number];
    margin?: number | [number, number] | [number, number, number, number];
    radius?: number | [number, number] | [number, number, number, number];
    borderWidth?: number | [number, number] | [number, number, number, number];
    gap?: number;
    pressedOpacity?: number;
    pressedColor?: BoxColor;
    children?: PressableProps['children'];
  };
const Box = ({
  pressedColor = 'textPrimary',
  pressedOpacity = 0.85,
  bgColor = 'surface',
  borderColor = 'divider',
  padding = 1,
  margin = 0,
  radius = 1,
  borderWidth = 1,
  gap = 0,
  children,
  onPress,
  onLongPress,
  onPressIn,
  onPressOut,
  style,
  ...props
}: BoxProps) => {
  const {colors, sizes} = useTheme();

  const pressabled = useMemo(
    () => !!(onPress || onLongPress || onPressIn || onPressOut),
    [onPress, onLongPress, onPressIn, onPressOut],
  );

  const backgroundColor_ = useMemo(
    () => transformColor(colors, bgColor),
    [bgColor, colors],
  );

  const borderColorStyle = useMemo<ViewStyle>(
    () => transformBorderColor(colors, borderColor),
    [borderColor, colors],
  );

  const paddingStyle = useMemo(
    () => transformPaddingStyle<ViewStyle>(sizes.base, padding),
    [padding, sizes],
  );

  const marginStyle = useMemo(
    () => transformMarginStyle<ViewStyle>(sizes.base, margin),
    [margin, sizes],
  );

  const radiusStyle = useMemo(
    () => transformRadiusStyle<ViewStyle>(sizes.base, radius),
    [radius, sizes],
  );

  const borderWidthStyle = useMemo(
    () => transformBorderWidthStyle<ViewStyle>(sizes.border, borderWidth),
    [borderWidth, sizes],
  );

  const gapStyle = useMemo(
    () => transformGapStyle<ViewStyle>(sizes.base, gap),
    [sizes, gap],
  );

  const style_ = useMemo(
    () =>
      StyleSheet.compose(style, {
        backgroundColor: backgroundColor_,
        ...borderColorStyle,
        ...paddingStyle,
        ...marginStyle,
        ...radiusStyle,
        ...borderWidthStyle,
        ...gapStyle,
      }),
    [
      backgroundColor_,
      borderColorStyle,
      paddingStyle,
      marginStyle,
      radiusStyle,
      borderWidthStyle,
      gapStyle,
    ],
  );

  const pressedColor_ = useMemo(
    () => (pressabled ? transformColor(colors, pressedColor) : void 0),
    [pressabled, pressedColor, colors],
  );

  const overrideStyle = useMemo(() => {
    if (!pressabled) {
      return void 0;
    }

    const overrideStyle: ViewStyle = {
      opacity: pressedOpacity,
    };

    for (const key of Object.keys(borderColorStyle)) {
      overrideStyle[
        key as
          | 'borderColor'
          | 'borderTopColor'
          | 'borderRightColor'
          | 'borderBottomColor'
          | 'borderLeftColor'
      ] = backgroundColor_;
    }
    return overrideStyle;
  }, [pressabled, backgroundColor_, borderColorStyle, pressedOpacity]);

  const pressableStyle = useMemo(
    () => ({
      backgroundColor: pressedColor_,
      ...radiusStyle,
    }),
    [pressedColor_, radiusStyle],
  );

  if (!pressabled) {
    return (
      <View style={style_} {...props}>
        {typeof children === 'function' ? children({pressed: false}) : children}
      </View>
    );
  }

  return (
    <Pressable
      style={({pressed}) => (pressed ? pressableStyle : void 0)}
      onPress={onPress}
      onLongPress={onLongPress}
      onPressIn={onPressIn}
      onPressOut={onPressOut}>
      {({pressed}) => (
        <View
          style={pressed ? StyleSheet.compose(style_, overrideStyle) : style_}
          {...props}>
          {typeof children === 'function' ? children({pressed}) : children}
        </View>
      )}
    </Pressable>
  );
};

export default Box;
