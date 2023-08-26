import React, { useEffect, useState } from 'react';
import { PingTest, PongTest, BroadcastName, createSocketClient, SocketContext, useBoston } from './socketsdk';


const client = createSocketClient("ws://localhost:7000/ws")


function Compo(){
  const [msg, setMsg] = useState<string>("")

  const {listen, send} = useBoston()

  
  useEffect(() => {
    const clear = listen('broadcast_name', data => {
      setMsg(data.name)
    })

    return clear
  }, [])


  useEffect(() => {
    const inter = setInterval(() => {
      send("ping_test", {
        data: "slow"
      })
    }, 3000)
    return () => {
      clearInterval(inter)
    }
  }, [])

  return (
    <div className="App">
      Broadcast all : {msg}
    </div>
  )
}


function App() {

  
  
  
  return (
    <SocketContext.Provider value={client}>
       <Compo />
    </SocketContext.Provider>
   
  );
}

export default App;
