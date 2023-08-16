export interface Namess {
	name: string
	da: string
	example: string
}

export interface DataT {
	datal: Array<number>
}

export interface ResData {
	names: Namess
	page: number
	data: DataT
	slite: Array<number>
	datap: DataT | undefined
	das: DataT | undefined
	dasss: Array<DataT | undefined>
	dmap: {[key: string]: string}
	dmapo: {[key: string]: DataT}
	dmap1: {[key: string]: string}
	dmapo1: {[key: string]: DataT}
}

export interface Access {
	access_type: string
}

export interface OtherInfo {
	whatsapp: string
	viagra: string | undefined
	access_: Access
}

export interface Time {
}

export interface Example {
	other: OtherInfo
	name: string
	age: number | undefined
	adult: boolean
	js: {[key: string]: OtherInfo | undefined}
	ks: Array<string>
	is: Array<string | undefined>
	date: Time
}



export const client = {
	GetTestdata: {
	query: undefined,
	response: undefined,
	body: {
			datal: [
					0
				]
		},
	method: 'GET' as const,
	url: 'testdata' as const
},
	PostTestdata23: {
	query: undefined,
	response: {
				other: {
					whatsapp: ``,
					viagra: `` as string | undefined,
					access_: {
						access_type: ``
					}
				},
				name: ``,
				age: 0 as number | undefined,
				adult: false,
				js: {},
				ks: [
						``
					],
				is: [
						`` as string | undefined
					],
				date: {

}
			} as Example | undefined,
	body: {
			datal: [
					0
				]
		},
	method: 'POST' as const,
	url: 'testdata23' as const
},
	PostTestdata33: {
	query: undefined,
	response: [
				{
					names: {
						name: ``,
						da: ``,
						example: ``
					},
					page: 0,
					data: {
						datal: [
								0
							]
					},
					slite: [
							0
						],
					datap: {
							datal: [
									0
								]
						} as DataT | undefined,
					das: {
							datal: [
									0
								]
						} as DataT | undefined,
					dasss: [
							{
								datal: [
										0
									]
							} as DataT | undefined
						],
					dmap: {},
					dmapo: {},
					dmap1: {},
					dmapo1: {}
				} as ResData | undefined
			],
	body: {
			datal: [
					0
				]
		},
	method: 'POST' as const,
	url: 'testdata33' as const
}
}