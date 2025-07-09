import {
  CommonActions,
  DefaultTheme,
  NavigationContainer,
  useNavigation,
  type Theme as NavigationTheme,
} from '@react-navigation/native';
import {createNativeStackNavigator} from '@react-navigation/native-stack';

import type {Theme} from '../Common/Theme';

export {CommonActions, NavigationContainer, useNavigation};

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
      card: colors.surface,
      text: colors.textPrimary,
      border: colors.divider,
      notification: colors.error,
    },
    fonts: DefaultTheme.fonts,
  };
};
