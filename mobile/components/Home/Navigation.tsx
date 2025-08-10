import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import {useNavigation} from '@react-navigation/native';
import type {ComponentProps} from 'react';
import {Suspense} from 'react';
import {ScreenLoading} from '../Common/ScreenBase';

export {useNavigation};

const Tab = createBottomTabNavigator();

export const Navigator = Tab.Navigator;
export const Screen = Tab.Screen;
export const Group = Tab.Group;

type NavigatorScreenLayout = ComponentProps<typeof Navigator>['screenLayout'];
type ScreenLayoutProps = Parameters<
  Extract<NavigatorScreenLayout, Function>
>[0];
export const ScreenLayout = ({children}: ScreenLayoutProps) => (
  <Suspense fallback={<ScreenLoading />}>{children}</Suspense>
);
