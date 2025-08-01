import {useMemo} from 'react';
import type {ViewProps, ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';
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

export type PaperProps = ViewProps & {
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
};

const Paper = ({
  bgColor = 'surface',
  borderColor = 'transparent',
  padding = 0,
  margin = 0,
  radius = 1,
  borderWidth = 1,
  gap = 0,
  children,
  style,
  ...props
}: PaperProps) => {
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
      style,
    ],
  );

  return (
    <View style={style_} {...props}>
      {children}
    </View>
  );
};

export default Paper;
