import {useMemo} from 'react';
import type {TextProps as NativeTextProps, TextStyle} from 'react-native';
import {Text as NativeText, StyleSheet} from 'react-native';
import {scale} from './SizeMatters';
import {
  transformBorderColor,
  transformBorderWidthStyle,
  transformMarginStyle,
  transformPaddingStyle,
  transformRadiusStyle,
} from './StyleSheet';
import type {Theme, ThemeColor} from './Theme';
import {transformColor, useTheme} from './Theme';

const TextSizeFactors = {
  text: 0,
  xsmall: -0.5,
  small: -0.25,
  title: 0.25,
};
type TextSizeKey = keyof typeof TextSizeFactors;
type ThemeFontsKey = keyof Theme['fonts'];
type TextProps = NativeTextProps & {
  color?: ThemeColor;
  size?: TextSizeKey;
  font?: ThemeFontsKey;
  bgColor?: ThemeColor;
  borderColor?:
    | ThemeColor
    | [ThemeColor, ThemeColor]
    | [ThemeColor, ThemeColor, ThemeColor, ThemeColor];
  padding?: number | [number, number] | [number, number, number, number];
  margin?: number | [number, number] | [number, number, number, number];
  radius?: number | [number, number] | [number, number, number, number];
  borderWidth?: number | [number, number] | [number, number, number, number];
};
const Text = ({
  bgColor = 'transparent',
  borderColor = 'transparent',
  padding = 0,
  margin = 0,
  radius = 0,
  borderWidth = 0,
  color = 'inherit',
  size = 'text',
  font = 'regular',
  children,
  style,
  ...props
}: TextProps) => {
  const {sizes, colors, fonts} = useTheme();

  const color_ = useMemo(() => transformColor(colors, color), [color, colors]);

  const size_ = useMemo(() => {
    const factor = TextSizeFactors[size];
    if (factor === void 0 || factor === 0) {
      return scale(sizes.text);
    }
    return scale(sizes.text + sizes.base * factor);
  }, [size, sizes]);

  const font_ = useMemo(() => {
    const fontStyle = fonts[font];
    if (fontStyle === void 0) {
      return fonts.regular;
    }
    return fontStyle;
  }, [font, fonts]);

  const backgroundColor_ = useMemo(
    () => transformColor(colors, bgColor),
    [bgColor, colors],
  );

  const borderColorStyle = useMemo<TextStyle>(
    () => transformBorderColor(colors, borderColor),
    [borderColor, colors],
  );

  const paddingStyle = useMemo(
    () => transformPaddingStyle<TextStyle>(sizes.base, padding),
    [padding, sizes],
  );

  const marginStyle = useMemo(
    () => transformMarginStyle<TextStyle>(sizes.base, margin),
    [margin, sizes],
  );

  const radiusStyle = useMemo(
    () => transformRadiusStyle<TextStyle>(sizes.base, radius),
    [radius, sizes],
  );

  const borderWidthStyle = useMemo(
    () => transformBorderWidthStyle<TextStyle>(sizes.border, borderWidth),
    [borderWidth, sizes],
  );

  const style_ = useMemo(
    () =>
      StyleSheet.compose(style, {
        color: color_,
        fontSize: size_,
        fontFamily: font_.fontFamily,
        fontWeight: font_.fontWeight as TextStyle['fontWeight'],
        backgroundColor: backgroundColor_,
        ...borderColorStyle,
        ...paddingStyle,
        ...marginStyle,
        ...radiusStyle,
        ...borderWidthStyle,
      }),
    [
      color_,
      size_,
      font_,
      backgroundColor_,
      borderColorStyle,
      paddingStyle,
      marginStyle,
      radiusStyle,
      borderWidthStyle,
      style,
    ],
  );

  return (
    <NativeText style={style_} {...props}>
      {children}
    </NativeText>
  );
};
export default Text;
