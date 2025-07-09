import {ComponentProps} from 'react';
import type {ViewProps, ViewStyle} from 'react-native';
import {ScrollView, StyleSheet, View} from 'react-native';
import {scale, verticalScale} from './SizeMatters';
import {withStyle, withTheme} from './StyleSheet';
import type {Theme} from './Theme';
import {useTheme} from './Theme';

const styles = StyleSheet.create({
  row: {
    flexDirection: 'row',
    flexWrap: 'wrap',
  },
  actions: {
    flexDirection: 'row',
    alignItems: 'flex-end',
  },
});

type LayoutProps = ComponentProps<typeof ScrollView>;
export const Layout = ({children, style, ...props}: LayoutProps) => {
  const {sizes} = useTheme();
  return (
    <ScrollView
      {...props}
      style={StyleSheet.compose({paddingHorizontal: scale(sizes.base)}, style)}>
      <View style={{height: verticalScale(sizes.base * 1.5)}} />
      {children}
    </ScrollView>
  );
};

export const RowLayout = withStyle<ViewStyle, ViewProps>(View, styles.row);
export const ActionsLayout = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes}) => ({
    ...styles.actions,
    gap: scale(sizes.base),
  }),
);
