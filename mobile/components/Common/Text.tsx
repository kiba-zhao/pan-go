import {useMemo} from 'react';
import type {
  ColorValue,
  TextProps as NativeTextProps,
  TextStyle,
} from 'react-native';
import {Text as NativeText, StyleSheet} from 'react-native';
import {scale} from './SizeMatters';
import {type Theme, useTheme} from './Theme';

const TextSizeFactors = {
  text: 0,
  xsmall: -0.5,
  small: -0.25,
  title: 0.25,
};
type TextSizeKey = keyof typeof TextSizeFactors;
type ThemeColorsKey = keyof Theme['colors'];
type ThemeFontsKey = keyof Theme['fonts'];
type TextProps = NativeTextProps & {
  color?: ThemeColorsKey | ColorValue;
  size?: TextSizeKey;
  font?: ThemeFontsKey;
};
const Text = ({
  color = 'inherit',
  size = 'text',
  font = 'regular',
  children,
  style,
  ...props
}: TextProps) => {
  const {sizes, colors, fonts} = useTheme();

  const color_ = useMemo(
    () =>
      typeof color === 'string' && colors[color as ThemeColorsKey] !== void 0
        ? colors[color as ThemeColorsKey]
        : color,
    [color, colors],
  );

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

  return (
    <NativeText
      style={StyleSheet.compose(
        {
          color: color_,
          fontSize: size_,
          fontFamily: font_.fontFamily,
          fontWeight: font_.fontWeight as TextStyle['fontWeight'],
          verticalAlign: 'middle',
        },
        style,
      )}
      {...props}>
      {children}
    </NativeText>
  );
};
export default Text;
