import {createBottomTabNavigator} from '@react-navigation/bottom-tabs';
import {NavigationContainer} from '@react-navigation/native';
import type {ComponentProps} from 'react';

const Tab = createBottomTabNavigator();

export type NavigatorProps = ComponentProps<typeof Tab.Navigator>;
export const Navigator = ({children, ...props}: NavigatorProps) => {
  return (
    <NavigationContainer>
      <Tab.Navigator {...props}>{children}</Tab.Navigator>
    </NavigationContainer>
  );
};

export type ScreenProps = ComponentProps<typeof Tab.Screen>;
export const Screen = Tab.Screen;

export type ScreenPropsOptions = ScreenProps['options'];
export type TabBarIconType = NonNullable<
  NonNullable<Exclude<ScreenPropsOptions, Function>>['tabBarIcon']
>;
