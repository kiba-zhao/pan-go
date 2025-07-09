import {ComponentType} from 'react';
import {
  StyleSheet,
  type ImageStyle,
  type StyleProp,
  type TextStyle,
  type ViewStyle,
} from 'react-native';
import {withTheme as withThemeBase, type Theme} from './Theme';

type CreatePropsFunc<
  BaseProps extends {},
  Props extends {},
  T extends Theme,
> = Parameters<typeof withThemeBase<BaseProps, Props, T>>[1];

function createStyleWithTheme<
  T extends Theme,
  Style extends ViewStyle | TextStyle | ImageStyle,
  Props extends {style?: StyleProp<Style>},
>(
  createStyleHandleFunc: (theme: T) => Style,
): CreatePropsFunc<Props, Props, T> {
  return (theme, props) => {
    const themeStyle = createStyleHandleFunc(theme);
    return {...props, style: StyleSheet.compose(themeStyle, props.style)};
  };
}

export function withTheme<
  T extends ViewStyle | TextStyle | ImageStyle,
  V extends {style?: StyleProp<T>},
  U extends Theme,
>(
  BaseComponent: ComponentType<V>,
  createStyle: (theme: U) => T,
): ComponentType<V> {
  return withThemeBase(BaseComponent, createStyleWithTheme(createStyle));
}

export function withStyle<
  T extends ViewStyle | TextStyle | ImageStyle,
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
  T extends ViewStyle | TextStyle | ImageStyle,
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
