
import { useEffect } from 'react';
import Header from './components/Header';
import RootLayout from './components/RootLayout';
import { usePlaygroundStore } from './lib/playgroundStore';

function App() {
    const [initPlayground] = usePlaygroundStore((state) => [state.initPlayground]);

    useEffect(() => {
        void initPlayground();
    }, [initPlayground])

    return <>
        <Header />
        <RootLayout />
    </>
}

export default App
