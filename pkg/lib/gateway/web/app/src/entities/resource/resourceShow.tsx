import React, {useCallback, useContext, useMemo, useState} from 'react'
import {IResourceComponentsProps, useShow, useUpdate} from '@refinedev/core'
import {Show, TextField} from '@refinedev/antd'
import {Button, Input, Typography} from 'antd'
import ReactJson from 'react-json-view'
import {ColorModeContext} from '../../contexts/color-mode'
import {Status} from './resourceList'

const {TextArea} = Input
const {Title} = Typography

export const ResourceShow: React.FC<IResourceComponentsProps> = () => {
    const {queryResult} = useShow()
    const {data, isLoading} = queryResult
    const record = data?.data || {}

    const [resolvedData, setResolvedData] = useState('')
    const {mutate} = useUpdate()
    const {mode} = useContext(ColorModeContext)

    const valueAsJson = useMemo(() => {
        if (record?.status !== 'ready') return undefined

        try {
            return {data: JSON.parse(record.value)}
        } catch {
            return undefined
        }
    }, [record])

    const handleResolveClick = useCallback(() => {
        if (!record?.id) return

        mutate({
            resource: 'resources',
            id: record.id,
            values: {value: resolvedData},
        })
    }, [record, resolvedData, mutate])

    const handleTextAreaChange = useCallback((event: React.ChangeEvent<HTMLTextAreaElement>) => {
        setResolvedData(event.target.value)
    }, [])

    return (
        <Show isLoading={isLoading}>
            <Title level={5}>ID</Title>
            <TextField copyable value={record?.id}/>

            <Title level={5}>Status</Title>
            <TextField value={<Status value={record?.status}/>}/>

            <Title level={5}>Data</Title>
            {record?.status === 'ready' && valueAsJson ? (
                <ReactJson
                    src={valueAsJson}
                    name={null}
                    theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
                    collapseStringsAfterLength={80}
                />
            ) : (
                <TextArea
                    value={record?.status === 'ready' ? record?.value : resolvedData}
                    placeholder='Put your data here'
                    disabled={record?.status === 'ready'}
                    onChange={handleTextAreaChange}
                    allowClear
                />
            )}

            <Button
                disabled={record?.status === 'ready'}
                onClick={handleResolveClick}
                style={{marginTop: 16}}
            >
                Resolve
            </Button>
        </Show>
    )
}
