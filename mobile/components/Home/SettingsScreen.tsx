import {AppSettings} from '@pango/data';
import {Fragment} from 'react';
import {CommonActions, useNavigation} from '../App/Navigation';
import {ScreenLayout} from '../App/ScreenBase';
import {IconFieldItem, TextFieldItem} from '../Common/Field';
import {useTranslation} from '../Common/I18Next';
import {useQuery} from '../Common/ReactQuery';
import {ScreenRefreshControl} from '../Common/ScreenBase';
import Section, {SectionHeader, SectionItem} from '../Common/Section';
import Text from '../Common/Text';
import {QueryKey as SettingsQueryKey} from '../Settings/ReactQuery';
import {
  BroadcastAddressEditorScreenName,
  NameEditorScreenName,
  PeerIDScreenName,
  PeerPortEditorScreenName,
  QRCodeScreenName,
} from '../Settings/ScreenRoute';
import {load as loadSettings} from '../Spec/AppSettings';
import {HomeDeviceScreenName} from './ScreenRoute';

const SettingsScreen = () => {
  const navigation = useNavigation();
  const handleSwipeRight = () => {
    navigation.dispatch(CommonActions.navigate(HomeDeviceScreenName));
  };

  const {refetch, data} = useQuery({
    queryKey: SettingsQueryKey,
    queryFn: loadSettings,
    placeholderData: _ => _,
  });

  const handleRefresh = () => {
    refetch();
  };

  return (
    <ScreenLayout
      style={{paddingHorizontal: 0}}
      refreshControl={<ScreenRefreshControl onRefresh={handleRefresh} />}
      swipeOpts={{
        onSwipeRight: handleSwipeRight,
      }}>
      <BaseSection settings={data} />
      <NetSection settings={data} />
    </ScreenLayout>
  );
};

export default SettingsScreen;

type BaseSectionProps = {settings?: AppSettings};
const BaseSection = ({settings}: BaseSectionProps) => {
  const {t} = useTranslation();
  const navigation = useNavigation();

  const handleQRCodePress = () => {
    navigation.dispatch(CommonActions.navigate(QRCodeScreenName));
  };

  const handleNamePress = () => {
    navigation.dispatch(CommonActions.navigate(NameEditorScreenName));
  };

  const handlePeerIDPress = () => {
    navigation.dispatch(CommonActions.navigate(PeerIDScreenName));
  };

  return (
    <Fragment>
      <Section>
        <SectionItem variant="row" onPress={handleQRCodePress}>
          <IconFieldItem
            label={t('screen.settings.sections.qrcode')}
            disabled={!settings}
            name="qr-code-outline"
          />
        </SectionItem>
      </Section>
      <Section>
        <SectionItem variant="row" onPress={handleNamePress}>
          <TextFieldItem
            label={t('screen.settings.sections.name')}
            text={settings?.name}
            editable
          />
        </SectionItem>
        <SectionItem variant="row" onPress={handlePeerIDPress}>
          <IconFieldItem
            label={t('screen.settings.sections.peerId')}
            disabled={!settings}
            name="open-outline"
          />
        </SectionItem>
      </Section>
    </Fragment>
  );
};

type NetSectionProps = BaseSectionProps;
const NetSection = ({settings}: NetSectionProps) => {
  const {t} = useTranslation();
  const navigation = useNavigation();

  const handlePeerPortPress = () => {
    navigation.dispatch(CommonActions.navigate(PeerPortEditorScreenName));
  };

  const handleBroadcastAddressPress = () => {
    navigation.dispatch(
      CommonActions.navigate(BroadcastAddressEditorScreenName),
    );
  };
  return (
    <Section>
      <SectionHeader>
        <Text size="small" color="textSecondary">
          {t('screen.settings.sections.netSection')}
        </Text>
      </SectionHeader>
      <SectionItem variant="row" onPress={handlePeerPortPress}>
        <TextFieldItem
          label={t('screen.settings.sections.peerPort')}
          text={settings?.peerPort ? settings.peerPort.toString() : ''}
          editable
        />
      </SectionItem>
      <SectionItem variant="row" onPress={handleBroadcastAddressPress}>
        <IconFieldItem
          label={t('screen.settings.sections.broadcastAddress')}
          disabled={
            !settings ||
            !settings.broadcastAddress ||
            settings.broadcastAddress.length <= 0
          }
          name="chevron-forward-outline"
        />
      </SectionItem>
      <SectionItem variant="row">
        <IconFieldItem
          label={t('screen.settings.sections.publicAddress')}
          disabled={
            !settings ||
            !settings.publicAddress ||
            settings.publicAddress.length <= 0
          }
          name="chevron-forward-outline"
        />
      </SectionItem>
    </Section>
  );
};
