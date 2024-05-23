import {Configuration, DefaultApi} from "../api/generated";

const API = new DefaultApi(new Configuration(), 'http://localhost:8080');

interface ConvertedFilter<T> {
  [key: string]: {
    Type: string;
    Data: T;
  };
}

function ConvertFilter<T>(filters: Record<string, T>[]): ConvertedFilter<T> {
  const convertedFilter: ConvertedFilter<T> = {};

  for (const filter of filters) {
    const { field, operator, value } = filter;
    if (field === 'requestId' || (Array.isArray(value) && value.length === 0)) {
        continue;
    }
    let convertedValue = value;
    if (typeof value === 'string' && !isNaN(Date.parse(value))) {
        convertedValue = (new Date(value)).getTime() as unknown as T;
    }

    convertedFilter[field] = {
      Type: operator,
      Data: convertedValue,
    };
  }

  return convertedFilter;
}

export {API, ConvertFilter}