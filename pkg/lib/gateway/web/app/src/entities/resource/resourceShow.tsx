import React, {useCallback, useContext, useState} from "react";
import {BaseRecord, IResourceComponentsProps, useShow, useUpdate, useUpdateMany} from "@refinedev/core";
import {Show, TagField, TextField} from "@refinedev/antd";
import {Button, Typography} from "antd";
import {Status} from "./resourceList";
import {Input} from 'antd';
import ReactJson from "react-json-view";
import {ColorModeContext} from "../../contexts/color-mode";

const {TextArea} = Input;
const {Title} = Typography;

export const ResourceShow: React.FC<IResourceComponentsProps> = () => {
    const {queryResult} = useShow();
    const {data, isLoading} = queryResult;

    const record = data?.data;

    const [resolvedData, setResolvedData] = useState('')

    const {mutate} = useUpdate()

    const {mode} = useContext(ColorModeContext);

    const valueAsJson = (() => {
        if (record?.status !== 'ready') {
            return undefined
        }
        try {
            return JSON.parse(record.value);
        } catch {
            return undefined;
        }
    })()

    return (
        <Show isLoading={isLoading}>
            <Title level={5}>ID</Title>
            <TextField copyable value={record?.id}/>
            <Title level={5}>Status</Title>
            <TextField value={<Status value={record?.status}/>}/>
            <Title level={5}>Data</Title>
            {record?.status === 'ready' && valueAsJson !== undefined ?
                <ReactJson
                    src={valueAsJson}
                    name={'data'}
                    theme={mode === 'light' ? 'summerfruit:inverted' : 'summerfruit'}
                    collapseStringsAfterLength={80}
                />
                :
                <TextArea
                    value={record?.status === 'ready' ? record?.value : resolvedData}
                    placeholder={"Put your data here"}
                    disabled={record?.status === 'ready'}
                    onChange={(event) => {
                        // console.log('value:', event.target.value)
                        setResolvedData(event.target.value)
                    }}
                    allowClear
                />
            }
            <br/>
            <br/>
            <Button
                disabled={record?.status === 'ready'}
                onClick={() => {
                    if (!record) return
                    if (!record.id) return

                    let response = mutate({
                        resource: 'resources',
                        id: record.id,
                        values: {
                            value: resolvedData
                        }
                    })

                    console.log(response)
                }}
            >
                Resolve
            </Button>

        </Show>
    );
};
