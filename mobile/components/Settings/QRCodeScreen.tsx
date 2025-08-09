import {AppSettings} from '@pango/data';
import {HeaderTitle, ScreenLayout} from '../App/ScreenBase';
import Paper from '../Common/Paper';
import {ScreenRefreshControl} from '../Common/ScreenBase';
import Text from '../Common/Text';
import {useSettings} from './ReactQuery';
import {DeviceQRCode, I18NextProvider} from './ScreenBase';

const QRCodeScreen = () => {
  const {data, refetch} = useSettings();

  return (
    <I18NextProvider>
      <ScreenLayout
        refreshControl={<ScreenRefreshControl onRefresh={refetch} />}>
        <HeaderTitle i18nKey="screen.qrcode.name" />
        <FieldsSection settings={data} />
      </ScreenLayout>
    </I18NextProvider>
  );
};

export default QRCodeScreen;

type FieldsSectionProps = {settings?: AppSettings};
const FieldsSection = ({settings}: FieldsSectionProps) => {
  return (
    <Paper
      bgColor="transparent"
      style={{alignItems: 'center'}}
      gap={1}
      padding={[2, 0, 0, 0]}>
      <DeviceQRCode
        size="medium"
        disabled={!settings}
        name={settings?.name || ''}
        peerId={settings?.peerId || ''}
      />
      <Text>{settings?.name || ''}</Text>
    </Paper>
  );
};
