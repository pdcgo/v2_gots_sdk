package v2_gots_sdk

import (
	"os"
	"strings"

	"github.com/tkrajina/typescriptify-golang-structs/typescriptify"
)

var eventSocketTemplate = `
	#EventName#: #DataName#`

var classScriptTemplate = `

type KeyEvent = keyof SocketSdkConfig

export interface EventMessage {
	event_key: KeyEvent
	body: string
}



type ListenData = {
	[key in KeyEvent]?: (event: any) => void
}

export class SdkWebsocket extends WebSocket {
	
	listenerdata: ListenData = {}

	constructor(url: string){
		super(url);
		this.onmessage = (event) => {
			const eventmessage: EventMessage = JSON.parse(event.data)
			const eventdata = atob(eventmessage.body)
			const eventp = JSON.parse(eventdata)
			
			const handler = this.listenerdata[eventmessage.event_key]
			if (handler !== undefined) {
				handler(eventp)
			}
			
		}
	}


	sdkSend<K extends keyof SocketSdkConfig>(key: K, data: SocketSdkConfig[K]){
		const encodedData = btoa(JSON.stringify(data))

		const event: EventMessage = {
			event_key: key,
			body: encodedData
		}

		return this.send(JSON.stringify(event)) 
	}
	
	setEventListener<K extends keyof SocketSdkConfig>(key: K, handler: (data: SocketSdkConfig[K]) => any){
		
		this.listenerdata[key] = handler
	}
}
`

func CreateTsSocketEvent(generator *typescriptify.TypeScriptify, event EventIface) string {
	template := eventSocketTemplate

	eventName := event.KeyEvent()
	dataname := getStructName(generator, event, false)

	template = strings.ReplaceAll(template, "#EventName#", eventName)
	template = strings.ReplaceAll(template, "#DataName#", dataname)

	return template
}

func (sdk *SocketGenerator) CreateRootTypeSocket(f *os.File, funcscripts []string) {
	f.WriteString("/* eslint-disable @typescript-eslint/no-explicit-any */\n\n")

	model, _ := sdk.Model.Convert(map[string]string{})
	f.WriteString(model)
	f.WriteString("\n")
	f.WriteString("export type SocketSdkConfig = { \n")
	f.WriteString(strings.Join(funcscripts, ",\n"))
	f.WriteString("\n}\n")
	f.WriteString(classScriptTemplate)
}
