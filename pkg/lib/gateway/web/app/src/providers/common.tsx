import {Configuration, DefaultApi} from "../api/generated";

const API = new DefaultApi(new Configuration(), 'http://localhost:8585');

interface ConvertedFilter<T> {
    [key: string]: {
        Type: string;
        Data: T;
    };
}

function filtersToMongoFilter(filters) {
    const result: { $and: any[] } = {$and: []}

    for (const filter of filters) {
        const {field, operator, value} = filter;
        if (field === 'requestId' || (Array.isArray(value) && value.length === 0)) {
            continue;
        }
        let convertedValue = value;
        if (typeof value === 'string' && !isNaN(Date.parse(value))) {
            convertedValue = Math.round((new Date(value)).getTime() / 1000) as unknown as T;
        }

        result.$and.push({
            [field]: {
                [operatorToMqlOperator(operator)]: convertedValue,
            }
        })
    }

    if (result.$and.length === 0) delete result['$and'];

    return result;
}

function operatorToMqlOperator(operator) {
    switch (operator) {
        case "lte":
            return "$lte"
        case "gte":
            return "$gte"
        case "in":
            return "$in"
        default:
            console.log('unknown operator:', operator)
            return ""
    }
}

export {API, filtersToMongoFilter, operatorToMqlOperator}