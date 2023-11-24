import React from 'react';
import logo from './logo.svg';
import './App.css';
import ChatJoin from './Join';
import Chat from './Chat';
import Navbar from './Navbar';
import AppRoot from './AppRoot';
import { useMantineColorScheme, Button, Group } from '@mantine/core';
import { IconSun, IconMoon } from '@tabler/icons-react';

function App() {

  return (
    <div className="App">
      {/* <ChatJoin /> */}
      {/* <Chat /> */}
      {/* <Navbar /> */}
      <AppRoot />
    </div>
  );
}

export default App;
