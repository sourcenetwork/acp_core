
import { useEffect } from 'react';
import Header from './components/Header';
import RootLayout from './components/RootLayout';
import { usePlaygroundStore } from './lib/acpHandler';

function App() {
    const { initialize } = usePlaygroundStore();

    useEffect(() => {
        initialize();
    }, [initialize])

    return <>
        <Header />
        <RootLayout />
    </>
}

export default App
