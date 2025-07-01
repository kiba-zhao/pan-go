import type {ComponentType, PropsWithChildren} from 'react';
import {createContext, useContext} from 'react';

import {withContext} from './Component';
import {darkColors, default as defaultColors} from './theming/colors';
import fonts from './theming/fonts';
import sizes from './theming/sizes';

export const DefaultTheme = {
  colors: defaultColors,
  fonts,
  sizes,
};

export const DarkTheme = {
  ...DefaultTheme,
  colors: darkColors,
};

export type Theme = {
  colors: typeof defaultColors | typeof darkColors;
  sizes: typeof sizes;
  fonts: typeof fonts;
};

const Context = createContext<Theme | null>(null);

export const ThemeProvider = <T extends Theme>({
  children,
  theme,
}: PropsWithChildren<{theme: T}>) => {
  return <Context.Provider value={theme}>{children}</Context.Provider>;
};

export const useTheme = <T extends Theme>() => useContext(Context) as T;

export const useThemeValue = <
  T extends Theme,
  K extends keyof T,
  V extends T[K],
>(
  key: K,
): V => useTheme<T>()[key] as V;

export function withTheme<
  BaseProps extends {},
  Props extends {},
  T extends Theme,
>(
  BaseComponent: ComponentType<BaseProps>,
  createProps: (theme: T, props: Props) => BaseProps,
): ComponentType<Props> {
  return withContext(BaseComponent, Context, createProps);
}
