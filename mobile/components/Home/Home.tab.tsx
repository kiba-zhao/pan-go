import type {ComponentProps} from 'react';

import {withTransform} from '../Common/Component';
import Icon from '../Common/Icon';

import {Screen} from './Home.navigation';

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
