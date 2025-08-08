import {HeaderButton} from '@react-navigation/elements';

import type {ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';

import type {ComponentProps} from 'react';
import {Fragment, useEffect} from 'react';
import {CommonActions, useNavigation, useRoute} from '../App/Navigation';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import {scale} from '../Common/SizeMatters';
import {withTheme} from '../Common/StyleSheet';
import {useTheme} from '../Common/Theme';

export const HeaderAction = withTheme(
  HeaderButton,
  ({sizes}) =>
    ({
      paddingVertical: scale(sizes.base),
      paddingHorizontal: scale(sizes.base),
    } as ViewStyle),
);

export const HeaderActions = withTheme(
  View,
  ({sizes}) =>
    ({
      alignItems: 'flex-end',
      flexDirection: 'row',
      marginRight: scale(sizes.base) * -1,
    } as ViewStyle),
);

type HeaderBackActionProps = Omit<
  ComponentProps<typeof HeaderButton>,
  'children'
> & {
  screenName?: string;
};
export const HeaderBackAction = ({
  screenName,
  onPress,
  style,
  ...props
}: HeaderBackActionProps) => {
  const {sizes} = useTheme();
  const route = useRoute();

  const navigation = useNavigation();
  const handlePress = () => {
    onPress?.();
    const action =
      screenName === void 0
        ? CommonActions.goBack()
        : CommonActions.navigate(screenName, {}, {pop: true});
    navigation.dispatch(action);
  };
  return (
    <Fragment>
      <HeaderButton
        style={StyleSheet.compose(
          {
            padding: scale(sizes.base),
            marginLeft: scale(sizes.base) * -1,
          },
          style,
        )}
        onPress={handlePress}
        {...props}>
        <Icon name="arrow-back-sharp" color="textPrimary" size="medium" />
      </HeaderButton>
      <View style={{width: scale(sizes.base * 2)}} />
    </Fragment>
  );
};

export const HeaderTitle = ({
  i18nKey,
  text,
}: {
  text?: string;
  i18nKey?: string;
}) => {
  const {t} = useTranslation();
  const navigation = useNavigation();
  useEffect(() => {
    navigation.setOptions({
      title: i18nKey ? t(i18nKey) : text,
    });
  }, []);
  return null;
};
