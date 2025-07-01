import {ComponentType} from 'react';
import {
  StyleSheet,
  type ImageStyle,
  type StyleProp,
  type TextStyle,
  type ViewStyle,
} from 'react-native';
import {useTheme, type Theme} from './Theme';

export function withTheme<
  T extends ViewStyle | TextStyle | ImageStyle,
  V extends {style?: StyleProp<T>},
  U extends Theme,
>(
  BaseComponent: ComponentType<V>,
  withThemeHandleFunc: (theme: U) => T,
): ComponentType<V> {
  return ({...props}: V) => {
    const theme = useTheme<U>();
    const themeStyle = withThemeHandleFunc(theme);
    return (
      <BaseComponent
        {...props}
        style={StyleSheet.compose(themeStyle, props.style)}
      />
    );
  };
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
