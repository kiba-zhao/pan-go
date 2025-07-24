import {
  CommonActions,
  DefaultTheme,
  NavigationContainer,
  StackActions,
  useNavigation,
  useRoute,
  type Theme as NavigationTheme,
  type StackActionType,
} from '@react-navigation/native';
import {
  createNativeStackNavigator,
  type NativeStackHeaderLeftProps,
} from '@react-navigation/native-stack';
import type {ComponentProps} from 'react';
import {Suspense} from 'react';

import {ScreenLoading} from '../Common/ScreenBase';
import type {Theme} from '../Common/Theme';

export {
  CommonActions,
  NavigationContainer,
  StackActions,
  useNavigation,
  useRoute,
};
export type {NativeStackHeaderLeftProps, StackActionType};

const Stack = createNativeStackNavigator();

export const Navigator = Stack.Navigator;
export const Screen = Stack.Screen;
export const Group = Stack.Group;

type NavigatorScreenLayout = ComponentProps<typeof Navigator>['screenLayout'];
type ScreenLayoutProps = Parameters<
  Extract<NavigatorScreenLayout, Function>
>[0];
export const ScreenLayout = ({children}: ScreenLayoutProps) => (
  <Suspense fallback={<ScreenLoading />}>{children}</Suspense>
);

export const createTheme = <T extends Theme>(
  theme: T,
  darkMode?: boolean,
): NavigationTheme => {
  const dark = !!darkMode;
  const {colors} = theme;
  return {
    dark,
    colors: {
      primary: colors.primary,
      background: colors.background,
      card: colors.surface,
      text: colors.textPrimary,
      border: colors.divider,
      notification: colors.error,
    },
    fonts: DefaultTheme.fonts,
  };
};
