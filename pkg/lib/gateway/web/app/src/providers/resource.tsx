import {DataProvider} from '@refinedev/core'
import {API} from './common'

export const ResourceProvider: DataProvider = {
    getList: async ({resource, pagination, sorters, filters}: {
                        resource?: string;
                        pagination?: { current?: number; pageSize?: number };
                        sorters?: { field: string; order: 'asc' | 'desc' }[];
                        filters?: Record<string, any>
                    } = {}
    ) => {
        const resource_ids = filters?.find(({field, value}: {
            field: string,
            value: string[]
        }) => field === 'id')?.value

        const {data: resources} = await API.resourcesLoadPost(resource_ids ?? [])

        return {
            data: resources,
            total: resources.length,
        }
    },

    getOne: async ({resource, id}: { resource: string; id: string }) => {
        const {data: resources} = await API.resourcesLoadPost([id])
        console.log('loading:', resources)
        return {
            data: resources[0],
        }
    },

    create: async ({resource, variables: {amount}}: { resource: string; variables: { amount: number } }) => {
        const {data: allocated} = await API.resourcesAllocatePost(amount)

        return {
            data: allocated,
        }
    },

    update: async ({resource, id, variables: {value}}: {
        resource: string;
        id: string;
        variables: { value: string }
    }) => {
        await API.resourcesInitializePost([{id: id, value: value}])
    },
}
