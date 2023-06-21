import React, { useEffect, useState } from 'react';
import { PingTest, PongTest, SdkWebsocket, BroadcastName } from './socketsdk';





function App() {

  const [count, setCount] = useState<number>(0)
  const [msg, setMsg] = useState<string>("")
  
  useEffect(()=>{
    

    const ws = new SdkWebsocket("ws://localhost:7000/ws")

    ws.setEventListener("pong_test", (event: PongTest) =>{
      console.log(event, event.data)
      setCount(event.data)
    })

    ws.setEventListener('broadcast_data', (event: BroadcastName) =>{
      setMsg(event.name)
    })

    const interval = setInterval(function(){
      ws.sdkSend('ping_test', {
        data: "test data"
      })
    }, 5000)

    return () => {
      console.log("teardown")
      ws.close()
      clearInterval(interval)
    }
  }, [])
  
  return (
    <div className="App">
      Test Count : {count}<br />
      Broadcast all : {msg}
    </div>
  );
}

export default App;
