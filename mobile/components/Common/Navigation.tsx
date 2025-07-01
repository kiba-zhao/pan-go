import {
  DefaultTheme,
  NavigationContainer,
  StackActions,
  useNavigation,
  type Theme as NavigationTheme,
} from '@react-navigation/native';
import {createNativeStackNavigator} from '@react-navigation/native-stack';

import type {Theme} from './Theme';

export {StackActions as Actions, NavigationContainer, useNavigation};

const Stack = createNativeStackNavigator();

export const Navigator = Stack.Navigator;
export const Screen = Stack.Screen;
export const Group = Stack.Group;

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
      card: colors.background,
      text: colors.onBackground,
      border: colors.outline,
      notification: colors.error,
    },
    fonts: DefaultTheme.fonts,
  };
};
