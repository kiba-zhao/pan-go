import type {ComponentProps} from 'react';

import {withTransform} from '../Common/Component';
import Icon from '../Common/Icon';

import {Screen} from './Navigation';

import {HeaderButton} from '@react-navigation/elements';
import type {ViewStyle} from 'react-native';
import {View} from 'react-native';

import {scale} from '../Common/SizeMatters';
import {withTheme} from '../Common/StyleSheet';

type ScreenPropsOptions = ComponentProps<typeof Screen>['options'];
type TabBarIcon = NonNullable<
  NonNullable<Exclude<ScreenPropsOptions, Function>>['tabBarIcon']
>;
type TabBarIconProps = ComponentProps<TabBarIcon>;

type IconProps = ComponentProps<typeof Icon>;
type TransformTabBarIconOptions = {
  focusedName: IconProps['name'];
  defaultName: IconProps['name'];
};
function transformTabBarIcon(
  {focused, size, color}: TabBarIconProps,
  {focusedName, defaultName}: TransformTabBarIconOptions,
) {
  return {
    name: focused ? focusedName : defaultName,
    size,
    color,
  };
}

export function newTabBarIcon({
  focusedName,
  defaultName,
}: TransformTabBarIconOptions) {
  return withTransform(Icon, transformTabBarIcon, {
    focusedName,
    defaultName,
  }) as TabBarIcon;
}

export const HeaderAction = HeaderButton;

export const HeaderActions = withTheme(
  View,
  ({sizes}) =>
    ({
      alignItems: 'flex-end',
      flexDirection: 'row',
      paddingRight: scale(sizes.base),
    } as ViewStyle),
);
