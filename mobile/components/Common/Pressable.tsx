import type {ComponentProps} from 'react';
import {useMemo} from 'react';
import type {PressableStateCallbackType, ViewStyle} from 'react-native';
import {Pressable as NativePressable, StyleSheet} from 'react-native';
import {
  transformBorderColor,
  transformBorderWidthStyle,
  transformGapStyle,
  transformMarginStyle,
  transformPaddingStyle,
  transformRadiusStyle,
} from './StyleSheet';
import type {ThemeColor} from './Theme';
import {transformColor, useTheme} from './Theme';

type NativePressableProps = ComponentProps<typeof NativePressable>;
export {PressableStateCallbackType};
export type PressableProps = {
  bgColor?: ThemeColor;
  borderColor?:
    | ThemeColor
    | [ThemeColor, ThemeColor]
    | [ThemeColor, ThemeColor, ThemeColor, ThemeColor];
  padding?: number | [number, number] | [number, number, number, number];
  margin?: number | [number, number] | [number, number, number, number];
  radius?: number | [number, number] | [number, number, number, number];
  borderWidth?: number | [number, number] | [number, number, number, number];
  gap?: number;
} & NativePressableProps;

const Pressable = ({
  bgColor = 'transparent',
  borderColor = 'transparent',
  padding = 0,
  margin = 0,
  radius = 0,
  borderWidth = 0,
  gap = 0,
  children,
  style,
  ...props
}: PressableProps) => {
  const {colors, sizes} = useTheme();

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

  const style_ = useMemo(() => {
    const customStyle = {
      backgroundColor: backgroundColor_,
      ...borderColorStyle,
      ...paddingStyle,
      ...marginStyle,
      ...radiusStyle,
      ...borderWidthStyle,
      ...gapStyle,
    } as ViewStyle;
    if (typeof style === 'function') {
      return (state: PressableStateCallbackType) =>
        StyleSheet.compose(style(state), customStyle);
    }
    return StyleSheet.compose(style, customStyle);
  }, [
    backgroundColor_,
    borderColorStyle,
    paddingStyle,
    marginStyle,
    radiusStyle,
    borderWidthStyle,
    gapStyle,
    style,
  ]);
  return (
    <NativePressable style={style_} {...props}>
      {children}
    </NativePressable>
  );
};

export default Pressable;
