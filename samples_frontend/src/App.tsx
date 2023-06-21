import React, { useEffect } from 'react';
import { PingTest, PongTest, SdkWebsocket } from './socketsdk';
import { json } from 'stream/consumers';


const ws = new SdkWebsocket("ws://localhost:7000/ws")

ws.setEventListener("pong_test", (event: PongTest) =>{
  console.log(event, event.data)
})


function App() {
  
  useEffect(function(){
    setInterval(function(){
      ws.sdkSend('ping_test', {
        data: "test data"
      })
    }, 1000)
  }, [])
  
  return (
    <div className="App">
      asdasdasdasd
    </div>
  );
}

export default App;
