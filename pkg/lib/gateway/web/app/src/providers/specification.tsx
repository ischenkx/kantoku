import {DataProvider} from "@refinedev/core";
import {API} from "./common";

export const SpecificationProvider: DataProvider = {
    getList: async ({resource}: {
                        resource?: string;
                    } = {}
    ) => {
        if (resource === 'specifications') {
            return await API.tasksSpecificationsGetAllPost();
        } else {
            return await API.tasksSpecificationsTypesGetAllPost();
        }
    },

    getOne: async ({resource, id}: { resource: string; id: string }) => {
        if (resource === 'specifications') {
            return await API.tasksSpecificationsGetPost({id: id});
        } else {
            return await API.tasksSpecificationsTypesGetPost({id: id});
        }
    },
};
