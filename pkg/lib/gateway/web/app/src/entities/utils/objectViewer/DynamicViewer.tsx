import React from 'react'
import {Collapse, GlobalToken, Space, theme, Typography} from 'antd'
import {CaretRightOutlined} from '@ant-design/icons'

const {Panel} = Collapse
const {Text, Paragraph} = Typography

const getPanelStyles = (token: GlobalToken) => ({
    // marginBottom: 24,
    // background: token.colorFillAlter,
    // borderRadius: token.borderRadiusLG,
    // border: 'none',
})

const StringViewer = ({label, data}: { label: string, data: string }) => (
    <Space direction='vertical' style={{width: '100%'}}>
        {label && <Text strong>{label}</Text>}
        <Text code ellipsis style={{maxWidth: '100%'}}>{data}</Text>
    </Space>
)

const NumberViewer = ({label, data}: { label: string, data: any }) => (
    <Space direction='vertical' style={{width: '100%'}}>
        {label && <Text strong>{label}</Text>}
        <Text code ellipsis style={{maxWidth: '100%'}}>{data}</Text>
    </Space>
)

const ArrayViewer = ({label, data, themeToken}: { label: string, data: any[], themeToken: GlobalToken }) => (
    <Space direction='vertical'>
        {label && <Text strong>{label}</Text>}
        {data && data.length > 0 ? data.map((item, index) => (
            <div style={{marginLeft: 16}}>
                <DynamicViewer key={index} label={`# ${index + 1}`} data={item} themeToken={themeToken}/>
            </div>
        )) : <Text>-</Text>}
    </Space>
)

const ObjectViewer = ({label, data, themeToken}: {
    label: string,
    data: Record<string, any>,
    themeToken: GlobalToken
}) => (
    <Collapse
        expandIcon={({isActive}) => <CaretRightOutlined rotate={isActive ? 90 : 0}/>}
    >
        <Panel header={label} key='1' style={{display: 'flex', flexDirection: 'column', ...getPanelStyles(themeToken)}}>
            <div style={{display: 'flex', flexDirection: 'column'}}>
                {Object.keys(data).map(subKey => (
                    <div style={{marginBottom: 8}}>
                        <DynamicViewer key={subKey} label={subKey} data={data[subKey]} themeToken={themeToken}/>
                    </div>
                ))}
            </div>
        </Panel>
    </Collapse>
)

const DynamicViewer = ({label, data, themeToken}: { label: string, data: any, themeToken: GlobalToken }) => {
   console.log('dv', data)

    if (Array.isArray(data)) {
        return <ArrayViewer label={label} data={data} themeToken={themeToken}/>
    } else if (data !== null && typeof data === 'object') {
        return <ObjectViewer label={label} data={data} themeToken={themeToken}/>
    } else if (typeof data === 'number') {
        return <NumberViewer label={label} data={data}/>
    } else if (typeof data === 'string') {
        return <StringViewer label={label} data={data}/>
    }

    console.log('no data')

    return <div>No data</div>
}

const Viewer = ({data, label}: { data: any, label: string }) => {
    const {token} = theme.useToken()
    return <DynamicViewer label={label} data={data} themeToken={token}/>
}

export default Viewer