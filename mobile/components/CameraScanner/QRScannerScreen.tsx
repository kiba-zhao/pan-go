import {UnmarshalDeviceQRCode} from '@pango/data';
import {useState} from 'react';
import {Alert, StyleSheet} from 'react-native';
import type {Code} from 'react-native-vision-camera';
import {Camera, useCodeScanner} from 'react-native-vision-camera';
import {CommonActions, useNavigation, useRoute} from '../App/Navigation.tsx';
import {DeviceCreatorScreenName} from '../Device/Screen.tsx';
import {CameraScannerScreen} from './ScreenBase.tsx';
import {type QRScannerParam, QRScannerScope} from './ScreenRoute.ts';

const QRScannerScreen = () => {
  const route = useRoute();
  const {scope} = route.params as QRScannerParam;

  const [active, setActive] = useState(true);

  const navigation = useNavigation();
  const codeScanner = useCodeScanner({
    codeTypes: ['qr'],
    onCodeScanned: codes => {
      if (!codes.length || codes.length < 0) return;
      setActive(false);

      let action = analysisDeviceAction(QRScannerScope.Device, codes);
      if (action !== void 0) {
        navigation.dispatch(action);
      } else {
        Alert.alert('识别错误', '未识别的二维码');
      }
      setActive(true);
    },
  });

  return (
    <CameraScannerScreen>
      {device => (
        <Camera
          style={StyleSheet.absoluteFill}
          device={device}
          isActive={active}
          codeScanner={codeScanner}
        />
      )}
    </CameraScannerScreen>
  );
};

export default QRScannerScreen;

function takeOne<T extends unknown>(
  codes: Code[],
  unmarshal: (code: string) => T | undefined,
): T | undefined {
  let results: T | undefined;
  for (const code of codes) {
    results = unmarshal(code.value || '');
    if (results !== void 0) break;
  }
  return results;
}

function analysisDeviceAction(
  scope: QRScannerScope,
  codes: Code[],
): CommonActions.Action | undefined {
  if (!scope?.includes(QRScannerScope.Device)) {
    return;
  }

  const deviceFields = takeOne(codes, UnmarshalDeviceQRCode);
  if (deviceFields === void 0) {
    return;
  }
  return CommonActions.navigate(DeviceCreatorScreenName, deviceFields);
}
