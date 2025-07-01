import {HeaderButton} from '@react-navigation/elements';
import {ComponentProps} from 'react';
import type {ViewProps, ViewStyle} from 'react-native';
import {View} from 'react-native';

import {scale} from 'react-native-size-matters';
import {withTheme} from '../Common/StyleSheet';
import {Theme} from '../Common/Theme';

export const HeaderAction = withTheme<
  ViewStyle,
  ComponentProps<typeof HeaderButton>,
  Theme
>(HeaderButton, ({sizes}) => ({
  padding: scale(sizes.base),
}));

export const HeaderActions = withTheme<ViewStyle, ViewProps, Theme>(
  View,
  ({sizes}) => ({
    flexDirection: 'row',
    alignItems: 'flex-end',
    gap: scale(sizes.base),
    marginRight: scale(sizes.base) * -1,
  }),
);
