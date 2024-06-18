import React from "react";
import {IResourceComponentsProps, useGo} from "@refinedev/core";
import {Create, useForm} from "@refinedev/antd";
import {Form, InputNumber} from "antd";


export const ResourceCreate: React.FC<IResourceComponentsProps> = () => {
    const formData = useForm({
        resource: 'resources',
        action: 'create',
        redirect: false,
        onMutationSuccess(data) {
            console.log('Allocated resources successfully:', data)
        },
        onMutationError(err) {
            console.log('Failed to allocate resources:', err)
        },
    });

    const {
        formProps,
        onFinish,
    } = formData

    const go = useGo();

    return (
        <Create saveButtonProps={
            {
                onClick() {
                    console.log('amount:', formProps.form?.getFieldValue('amount') || 1)

                    const data = {
                        amount: formProps.form?.getFieldValue('amount') || 1,
                    }

                    onFinish(data).then((result) => {
                        return go({
                            to: {
                                resource: "resources", // resource name or identifier
                                action: "list",
                            },
                            query: {
                                filters: [
                                    {
                                        field: "id",
                                        operator: "in",
                                        value: result?.data ?? [],
                                    },
                                ],
                            },
                            type: "replace",
                        })
                    })
                }
            }
        }>
            <Form {...formProps} layout="vertical">
                <Form.Item
                    label="Amount"
                    name="amount"
                    rules={[
                        {
                            required: true,
                        },
                    ]}
                >
                    <InputNumber min={1} defaultValue={1} value={1}/>
                </Form.Item>
            </Form>
        </Create>
    );
};
