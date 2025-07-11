import {useMemo} from 'react';
import type {ViewProps} from 'react-native';
import {StyleSheet, View} from 'react-native';
import {useTheme} from './Theme';

type DividerProps = {
  mode?: 'top' | 'bottom' | 'left' | 'right';
} & ViewProps;
export const Divider = ({
  mode = 'bottom',
  children,
  style,
  ...props
}: DividerProps) => {
  const {colors} = useTheme();

  const style_ = useMemo(() => {
    if (children === void 0)
      return {
        width: mode === 'top' || mode === 'bottom' ? 0 : 1,
        height: mode === 'left' || mode === 'right' ? 0 : 1,
        backgroundColor: colors.divider,
      };

    if (mode === 'top')
      return {borderTopWidth: 1, borderTopColor: colors.divider};
    if (mode === 'bottom')
      return {borderBottomWidth: 1, borderBottomColor: colors.divider};
    if (mode === 'left')
      return {borderLeftWidth: 1, borderLeftColor: colors.divider};
    if (mode === 'right')
      return {borderRightWidth: 1, borderRightColor: colors.divider};
    return void 0;
  }, [children, mode, colors]);

  return (
    <View style={StyleSheet.compose(style_, style)} {...props}>
      {children}
    </View>
  );
};
