import React, {Profiler} from 'react'
import {createRoot} from 'react-dom/client'

import App from './App'
import * as qs from 'qs'

import './style.css'

const qsParse = qs.parse

qs.parse = (str, opts) => {
    opts ||= {}
    opts.arrayLimit = 10000000
    return qsParse(str, opts)
}

const container = document.getElementById('root') as HTMLElement
const root = createRoot(container)

root.render(
    <React.StrictMode>
        <Profiler id={'App'} onRender={() => console.log('rendering the app')}>
            <App/>
        </Profiler>
    </React.StrictMode>
)
