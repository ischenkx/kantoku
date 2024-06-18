import {BaseRecord, DataProvider, GetListParams, GetListResponse} from "@refinedev/core";
import {TaskProvider} from "./task";
import {ResourceProvider} from "./resource";
import {SpecificationProvider} from "./specification";


export const ProviderRouter: DataProvider = {
    getList: async (params): Promise<GetListResponse> => {
        switch (params.resource) {
            case "tasks":
                return TaskProvider.getList(params);
            case "specifications":
            case "types":
                return SpecificationProvider.getList(params);
            case "resources":
                return ResourceProvider.getList(params);
            default:
                throw new Error(`Unknown resource: ${params.resource}`);
        }
    },
    getOne: async (params) => {
        switch (params.resource) {
            case "tasks":
                return TaskProvider.getOne(params);
            case "specifications":
            case "types":
                return SpecificationProvider.getOne(params);
            case "resources":
                return ResourceProvider.getOne(params);
            default:
                throw new Error(`Unknown resource: ${params.resource}`);
        }
    },

    create: async (params) => {
        console.log('CREATING', params)
        switch (params.resource) {
            case "tasks":
                return TaskProvider.create(params);
            case "specifications":
            case "types":
                return SpecificationProvider.create(params);
            case "resources":
                return ResourceProvider.create(params);
            default:
                throw new Error(`Unknown resource: ${resource}`);
        }
    },

    update: async (params) => {
        switch (params.resource) {
            case "tasks":
                return TaskProvider.update(params);
            case "specifications":
            case "types":
                return SpecificationProvider.update(params);
            case "resources":
                return ResourceProvider.update(params);
            default:
                throw new Error(`Unknown resource: ${params.resource}`);
        }
    }
};