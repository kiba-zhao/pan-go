import Clipboard from '@react-native-clipboard/clipboard';
import {
  HeaderAction,
  HeaderActions,
  HeaderTitle,
  ScreenLayout,
} from '../App/ScreenBase';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import {ScreenRefreshControl} from '../Common/ScreenBase';
import Text from '../Common/Text';
import {useSettings} from './ReactQuery';
import {I18NextProvider} from './ScreenBase';

const PeerIDScreen = () => {
  const {data, refetch} = useSettings();

  return (
    <I18NextProvider>
      <ScreenLayout
        refreshControl={<ScreenRefreshControl onRefresh={refetch} />}>
        <HeaderTitle i18nKey="screen.peerID.name" />
        {data?.peerId !== void 0 && data?.peerId.length > 0 ? (
          <PeerIDSection value={data?.peerId} />
        ) : (
          <PeerIDEmpty />
        )}
      </ScreenLayout>
    </I18NextProvider>
  );
};

export default PeerIDScreen;

export const PeerIDHeaderActions = () => {
  return (
    <HeaderActions>
      <PeerIDCopyAction />
    </HeaderActions>
  );
};

const PeerIDCopyAction = () => {
  const {data} = useSettings();

  const disabled = !data?.peerId || data?.peerId.length <= 0;
  const handlePress = () => {
    if (!data?.peerId) return;
    Clipboard.setString(data.peerId);
  };
  return (
    <HeaderAction onPress={handlePress} disabled={disabled}>
      <Icon
        name="open-outline"
        color={disabled ? 'textDisabled' : 'textPrimary'}
        size="medium"
      />
    </HeaderAction>
  );
};

const PeerIDEmpty = () => {
  const {t} = useTranslation();
  return (
    <Paper
      bgColor={'transparent'}
      padding={[8, 0]}
      style={{alignItems: 'center'}}>
      <Text font="bold" color="textDisabled" size="title">
        {t('screen.errors.empty')}
      </Text>
    </Paper>
  );
};

const PeerIDSection = ({value}: {value: string}) => {
  return (
    <Paper padding={1} borderColor="divider">
      <Text size="small" color="textSecondary" lines={5}>
        {value}
      </Text>
    </Paper>
  );
};
