import type {DeviceQRCodeValue} from '@pango/data';
import {MarshalDeviceQRCode} from '@pango/data';
import type {ComponentProps} from 'react';
import {useMemo} from 'react';
import QRCode from 'react-native-qrcode-svg';
import Button from '../Common/Button';
import Paper from '../Common/Paper';
import type {ScreenModalProps} from '../Common/ScreenBase';
import {ScreenModal} from '../Common/ScreenBase';
import {scale} from '../Common/SizeMatters';
import Text from '../Common/Text';
import {useTheme} from '../Common/Theme';

type QRCodeProps = ComponentProps<typeof QRCode>;
type DeviceQRCodeProps = DeviceQRCodeValue &
  Omit<QRCodeProps, 'value'> & {
    disabled?: boolean;
  };
export const DeviceQRCode = ({
  name,
  peerId,
  disabled,
  ...props
}: DeviceQRCodeProps) => {
  const {colors} = useTheme();
  const value = useMemo(
    () =>
      name.length > 0 && peerId.length > 0 && !disabled
        ? MarshalDeviceQRCode({name, peerId})
        : void 0,
    [disabled, name, peerId],
  );

  const color = useMemo(() => {
    if (disabled) return colors.textDisabled;
    return value === void 0 ? colors.error : colors.textPrimary;
  }, [disabled, value, colors]);

  return <QRCode size={scale(120)} color={color} value={value} {...props} />;
};

export const QRCodeSize = scale(280);

type DeviceQRCodeModalProps = Omit<ScreenModalProps, 'children'> &
  DeviceQRCodeValue;
export const DeviceQRCodeModal = ({
  name,
  peerId,
  onClose,
  ...props
}: DeviceQRCodeModalProps) => {
  return (
    <ScreenModal
      containerProps={{style: {alignItems: 'center'}}}
      onClose={onClose}
      {...props}>
      <DeviceQRCode size={QRCodeSize} name={name} peerId={peerId} />
      <Text size="small" margin={[1.5, 0, 0, 0]}>
        设备: {name}
      </Text>
      <Paper
        bgColor="transparent"
        padding={[3, 0, 1, 0]}
        gap={1.5}
        style={{width: QRCodeSize}}>
        <Button
          onPress={event => onClose?.(event)}
          title="关闭"
          disabled={name.length <= 0 || peerId.length <= 0}
        />
      </Paper>
    </ScreenModal>
  );
};
