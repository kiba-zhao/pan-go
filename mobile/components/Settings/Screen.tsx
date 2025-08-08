import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';
import {addLanguageSource} from '../Common/I18Next';
import {NameEditorScreenName, Namespace, QRCodeScreenName} from './ScreenRoute';

addLanguageSource('zh-CN', Namespace, () =>
  import(`../../locales/zh-CN/settings.json`).then(m => m.default),
);
addLanguageSource('en', Namespace, () =>
  import(`../../locales/en/settings.json`).then(m => m.default),
);

const SettingsAppScreen = () => (
  <Fragment>
    {QRCodeAppScreen()}
    {NameEditorAppScreen()}
  </Fragment>
);

export default SettingsAppScreen;

const QRCodeScreen = lazy(() => import('./QRCodeScreen'));

const QRCodeAppScreen = () => (
  <Screen
    name={QRCodeScreenName}
    component={QRCodeScreen}
    options={{title: ''}}
  />
);

const NameEditorScreen = lazy(() => import('./NameEditorScreen'));
const NameEditorAppScreen = () => (
  <Screen
    name={NameEditorScreenName}
    component={NameEditorScreen}
    options={{title: ''}}
  />
);
