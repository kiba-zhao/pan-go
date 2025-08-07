import {SafeAreaProvider} from 'react-native-safe-area-context';
import AppScreen from './components/App/Screen';
import {Provider as I18nextProvider} from './components/Common/I18Next';
import {ReactQueryProvider} from './components/Common/ReactQuery';

const App = () => (
  <SafeAreaProvider>
    <ReactQueryProvider>
      <I18nextProvider>
        <AppScreen />
      </I18nextProvider>
    </ReactQueryProvider>
  </SafeAreaProvider>
);

export default App;
