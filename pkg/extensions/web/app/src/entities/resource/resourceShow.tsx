import React, {useCallback, useState} from "react";
import {BaseRecord, IResourceComponentsProps, useShow, useUpdateMany} from "@refinedev/core";
import { Show, TagField, TextField } from "@refinedev/antd";
import {Button, Typography } from "antd";
import {Status} from "./resourceList";
import { Input } from 'antd';

const { TextArea } = Input;
const { Title } = Typography;

export const ResourceShow: React.FC<IResourceComponentsProps> = () => {
    const { queryResult } = useShow();
    const { data, isLoading } = queryResult;

    const record = data?.data;

    let [resolvedData, setResolvedData] = useState('')

    let {mutate} = useUpdateMany()


    return (
        <Show isLoading={isLoading}>
            <Title level={5}>ID</Title>
            <TextField copyable value={record?.id} />
            <Title level={5}>Status</Title>
            <TextField value={<Status value={record?.status}/>} />
            <Title level={5}>Data</Title>
            <TextArea
                value={record?.data}
                placeholder={record?.status === "ready" ? "Already resolved" : "Put your data here"}
                disabled={record?.status === 'ready'}
                onChange={(event)=>{
                    // console.log('value:', event.target.value)
                    setResolvedData(event.target.value)
                }}
                allowClear
            />
            <br/>
            <br/>
            <Button
                disabled={record?.status === 'ready'}
                onClick={()=>{
                    if (!record) return
                    if (!record.id) return

                    let response = mutate({
                        resource: 'resources',
                        ids: [record.id],
                        values: { data: resolvedData, status: 'ready' },
                    })

                    console.log(response)
                }}
            >
                Resolve
            </Button>

        </Show>
    );
};
