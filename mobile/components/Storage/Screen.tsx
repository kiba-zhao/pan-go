import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';
import {addLanguageSource} from '../Common/I18Next';
import {ExplorerScreenName, Namespace} from './ScreenRoute';

addLanguageSource('zh-CN', Namespace, () =>
  import(`../../locales/zh-CN/storage.json`).then(m => m.default),
);
addLanguageSource('en', Namespace, () =>
  import(`../../locales/en/storage.json`).then(m => m.default),
);

const StorageAppScreen = () => <Fragment>{ExplorerAppScreen()}</Fragment>;
export default StorageAppScreen;

const ExplorerScreen = lazy(() => import('./ExplorerScreen'));

const ExplorerAppScreen = () => (
  <Screen
    name={ExplorerScreenName}
    component={ExplorerScreen}
    options={{
      title: '',
    }}
  />
);
