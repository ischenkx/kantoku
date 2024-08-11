import React, {useEffect} from 'react'
import {IResourceComponentsProps, useDataProvider, useParsed} from '@refinedev/core'
import CytoscapeComponent from 'react-cytoscapejs'

import cytoscape from 'cytoscape'
import dagre from 'cytoscape-dagre'
import cose_bilknet from 'cytoscape-cose-bilkent'
import cise from 'cytoscape-cise'
import cola from 'cytoscape-cola'
import fcose from 'cytoscape-fcose'
// import elk from 'cytoscape-elk';

cytoscape.use(dagre)
cytoscape.use(cose_bilknet)
cytoscape.use(cise)
cytoscape.use(cola)
cytoscape.use(fcose)

// cytoscape.use( elk );

function getRandomColor() {
    // Generate random values for red, green, and blue components
    const red = Math.floor(Math.random() * 256)
    const green = Math.floor(Math.random() * 256)
    const blue = Math.floor(Math.random() * 256)

    // Create a CSS-formatted color string
    const color = `rgb(${red}, ${green}, ${blue})`

    return color
}

export const Sandbox: React.FC<IResourceComponentsProps> = () => {
    let {params} = useParsed<{ context: string }>()
    // let {data, isLoading, isError, refetch} = useList({
    //
    // })

    let dataProvider = useDataProvider()()

    let initialized = false
    let _cy = null

    useEffect(() => {
        const interval = setInterval(() => {
            dataProvider
                .getList({
                    resource: 'tasks',
                    filters: [
                        {
                            field: 'context',
                            operator: 'in',
                            value: [params.context],
                        }
                    ],
                    pagination: {
                        mode: 'off'
                    }
                })
                .then(
                    (result) => {
                        if (!_cy) return

                        let tasks = result.data || []

                        if (!initialized) {
                            _cy.$().remove()
                            console.log('initializing!')
                            let resource2task = {}
                            for (const task of tasks) {
                                for (const out of task.outputs) {
                                    resource2task[out] = task.id
                                }
                            }

                            for (const task of tasks) {
                                let color = 'blue'

                                _cy.add({
                                    data: {
                                        id: task.id,
                                        label: '',
                                    },
                                    style: {backgroundColor: color},
                                })
                            }

                            let edges = {}

                            for (const task of tasks) {
                                for (const input of task.inputs) {
                                    let responsibleTask = resource2task[input]
                                    if (!responsibleTask) continue
                                    let edgeId = `${responsibleTask}-${task.id}`
                                    if (edges[edgeId]) continue
                                    edges[edgeId] = true

                                    _cy.add({
                                        data: {
                                            id: edgeId,
                                            label: '',
                                            source: responsibleTask,
                                            target: task.id
                                        },
                                        // style: {
                                        //     'target-arrow-shape': 'triangle',
                                        //     'arrow-scale': 10,
                                        //     'target-arrow-color': 'red',
                                        // }
                                    })
                                }
                            }

                            _cy.makeLayout({name: 'dagre', padding: 15, align: 'DL'}).run()

                            console.log('initialized', getRandomColor(), _cy)
                            initialized = true
                        }

                        _cy.batch(() => {
                            for (let task of tasks) {
                                let color
                                switch (task.info.sub_status) {
                                    case 'ok':
                                        color = 'green'
                                        break
                                    case 'failed':
                                        color = 'red'
                                        break
                                    case 'initialized':
                                        color = 'blue'
                                        break
                                    default:
                                        color = 'yellow'
                                }

                                let currentColor = _cy.$(task.id).style()?.backgroundColor

                                if (currentColor === color) continue

                                // console.log(currentColor, '->', color)

                                // if (Math.random() < 0.5) continue

                                _cy.$id(task.id).style('backgroundColor', color)
                            }
                        })
                        console.log('updating...')
                    }
                )
        }, 1000)
        return () => {
            console.log('clearing...')
            clearInterval(interval)
        }
    }, [])

    const layoutOptions = [
        'cose',
        'cose-bilkent',
        'cise',
        'fcose',
        'cola',
        'spread',
        'concentric',
        'grid',
        'breadthfirst',
        'dagre',
        'klay',
        'random',
        'avsdf',
        'preset',
        'circle',
    ]

    const handleLayoutChange = (event) => {
        const selectedValue = event.target.value

        if (_cy) _cy.makeLayout({name: selectedValue}).run()
    }

    return (
        <>
            <div>
                <label htmlFor='layoutSelector'>Select Layout:</label>
                <select
                    id='layoutSelector'
                    onChange={handleLayoutChange}
                >
                    {layoutOptions.map((layout) => (
                        <option key={layout} value={layout}>
                            {layout}
                        </option>
                    ))}
                </select>
            </div>
            <CytoscapeComponent
                elements={[]}
                style={{width: '100%', height: '600px'}}
                cy={(cy) => {
                    initialized = false
                    _cy = cy

                    // for (let layout of layoutOptions) Cytoscape.use(layout)
                }}
            />
        </>

    )
}
