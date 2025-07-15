import {ComponentType} from 'react';
import {
  StyleSheet,
  type ColorValue,
  type ImageStyle,
  type StyleProp,
  type TextStyle,
  type ViewStyle,
} from 'react-native';
import {scale} from './SizeMatters';
import {transformColor, withTheme as withThemeBase, type Theme} from './Theme';

type CreatePropsFunc<
  BaseProps extends {},
  Props extends {},
  T extends Theme,
> = Parameters<typeof withThemeBase<BaseProps, Props, T>>[1];

function createStyleWithTheme<
  T extends Theme,
  Style extends ViewStyle | ImageStyle,
  Props extends {style?: StyleProp<Style>},
>(
  createStyleHandleFunc: (theme: T) => Style,
): CreatePropsFunc<Props, Props, T> {
  return (theme, {style, ...props}) => {
    const themeStyle = createStyleHandleFunc(theme);
    return {...(props as Props), style: StyleSheet.compose(themeStyle, style)};
  };
}

export function withTheme<
  T extends ViewStyle | ImageStyle,
  V extends {style?: StyleProp<T>},
  U extends Theme,
>(
  BaseComponent: ComponentType<V>,
  createStyle: (theme: U) => T,
): ComponentType<V> {
  return withThemeBase(BaseComponent, createStyleWithTheme(createStyle));
}

export function withStyle<
  T extends ViewStyle | ImageStyle,
  V extends {style?: StyleProp<T>},
>(BaseComponent: ComponentType<V>, baseStyle: StyleProp<T>) {
  return ({...props}: V) => {
    return (
      <BaseComponent
        {...props}
        style={StyleSheet.compose(baseStyle, props.style)}
      />
    );
  };
}

type StyleFunction<
  T extends ViewStyle | ImageStyle,
  Args extends Array<any>,
> = (...args: Args) => StyleProp<T>;

export function withStyleFunction<
  T extends ViewStyle | TextStyle | ImageStyle,
  Args extends Array<any>,
  V extends {style?: StyleProp<T> | StyleFunction<T, Args>},
>(
  BaseComponent: ComponentType<V>,
  baseStyle: StyleProp<T> | StyleFunction<T, Args>,
) {
  return ({...props}: V) => {
    const generateStyle = (...args: Args): StyleProp<T> => {
      const defaultStyle =
        typeof baseStyle === 'function' ? baseStyle(...args) : baseStyle;
      const customStyle =
        typeof props.style === 'function' ? props.style(...args) : props.style;
      return StyleSheet.compose(defaultStyle, customStyle);
    };
    return (
      <BaseComponent
        {...props}
        style={generateStyle as StyleFunction<T, Args>}
      />
    );
  };
}

export function transformPaddingStyle<Style extends ViewStyle | ImageStyle>(
  base: number,
  padding: number | [number, number] | [number, number, number, number],
): Style {
  if (typeof padding === 'number') {
    return {padding: scale(padding * base)} as Style;
  }
  if (Array.isArray(padding)) {
    if (padding.length === 4)
      return {
        paddingTop: scale(padding[0] * base),
        paddingRight: scale(padding[1] * base),
        paddingBottom: scale(padding[2] * base),
        paddingLeft: scale(padding[3] * base),
      } as Style;
    if (padding.length === 2)
      return {
        paddingHorizontal: scale(padding[1] * base),
        paddingVertical: scale(padding[0] * base),
      } as Style;
  }
  return {} as Style;
}

export function transformMarginStyle<Style extends ViewStyle | ImageStyle>(
  base: number,
  margin: number | [number, number] | [number, number, number, number],
): Style {
  if (typeof margin === 'number') {
    return {margin: scale(margin * base)} as Style;
  }
  if (Array.isArray(margin)) {
    if (margin.length === 4)
      return {
        marginTop: scale(margin[0] * base),
        marginRight: scale(margin[1] * base),
        marginBottom: scale(margin[2] * base),
        marginLeft: scale(margin[3] * base),
      } as Style;
    if (margin.length === 2)
      return {
        marginHorizontal: scale(margin[1] * base),
        marginVertical: scale(margin[0] * base),
      } as Style;
  }
  return {} as Style;
}

export function transformRadiusStyle<Style extends ViewStyle | ImageStyle>(
  base: number,
  radius: number | [number, number] | [number, number, number, number],
): Style {
  if (typeof radius === 'number') {
    return {borderRadius: scale(radius * base)} as Style;
  }
  if (Array.isArray(radius)) {
    if (radius.length === 4)
      return {
        borderTopLeftRadius: scale(radius[0] * base),
        borderTopRightRadius: scale(radius[1] * base),
        borderBottomRightRadius: scale(radius[2] * base),
        borderBottomLeftRadius: scale(radius[3] * base),
      } as Style;
    if (radius.length === 2)
      return {
        borderTopLeftRadius: scale(radius[0] * base),
        borderTopRightRadius: scale(radius[1] * base),
        borderBottomRightRadius: scale(radius[0] * base),
        borderBottomLeftRadius: scale(radius[1] * base),
      } as Style;
  }
  return {} as Style;
}

export function transformBorderWidthStyle<Style extends ViewStyle | ImageStyle>(
  base: number,
  width: number | [number, number] | [number, number, number, number],
): Style {
  if (typeof width === 'number') {
    return {borderWidth: scale(base * width)} as Style;
  }
  if (Array.isArray(width)) {
    if (width.length === 4)
      return {
        borderTopWidth: scale(base * width[0]),
        borderRightWidth: scale(base * width[1]),
        borderBottomWidth: scale(base * width[2]),
        borderLeftWidth: scale(base * width[3]),
      } as Style;
    if (width.length === 2)
      return {
        borderTopWidth: scale(base * width[0]),
        borderRightWidth: scale(base * width[1]),
        borderBottomWidth: scale(base * width[0]),
        borderLeftWidth: scale(base * width[1]),
      } as Style;
  }
  return {} as Style;
}

export function transformGapStyle<Style extends ViewStyle | ImageStyle>(
  base: number,
  gap: number | [number, number],
): Style {
  if (typeof gap === 'number') {
    return {gap: scale(gap * base)} as Style;
  }
  if (Array.isArray(gap)) {
    if (gap.length === 2)
      return {
        rowGap: scale(gap[0] * base),
        columnGap: scale(gap[1] * base),
      } as Style;
  }
  return {} as Style;
}

export function transformBorderColor<
  Colors extends Theme['colors'],
  Style extends ViewStyle | ImageStyle,
  Color extends NoInfer<keyof Colors> | ColorValue,
>(
  colors: Colors,
  color: Color | [Color, Color] | [Color, Color, Color, Color],
): Style {
  if (Array.isArray(color)) {
    if (color.length === 4)
      return {
        borderTopColor: transformColor(colors, color[0]),
        borderRightColor: transformColor(colors, color[1]),
        borderBottomColor: transformColor(colors, color[2]),
        borderLeftColor: transformColor(colors, color[3]),
      } as Style;
    if (color.length === 2) {
      const borderHorizontalColor = transformColor(colors, color[1]);
      const borderVerticalColor = transformColor(colors, color[0]);

      return {
        borderTopColor: borderVerticalColor,
        borderRightColor: borderHorizontalColor,
        borderBottomColor: borderVerticalColor,
        borderLeftColor: borderHorizontalColor,
      } as Style;
    }
    return {} as Style;
  }
  return {borderColor: transformColor(colors, color as Color)} as Style;
}
