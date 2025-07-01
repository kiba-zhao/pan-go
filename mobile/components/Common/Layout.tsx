import type {ViewProps, ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';
import {scale, verticalScale} from './SizeMatters';
import {withStyle, withTheme} from './StyleSheet';
import type {Theme} from './Theme';

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

export const Layout = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes}) => ({
    paddingTop: verticalScale(sizes.base * 1.5),
    paddingHorizontal: scale(sizes.base),
  }),
);

export const RowLayout = withStyle<ViewStyle, ViewProps>(View, styles.row);
export const ActionsLayout = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes}) => ({
    ...styles.actions,
    gap: scale(sizes.base),
  }),
);
