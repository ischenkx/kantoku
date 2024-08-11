export type Type = {
    name: string
    // TODO: use subTypes
    sub_types: Record<string, Type>
}

export type ResourceSet = {
    naming: {index: number, name: string}[]
    types: {index: number, type: Type}[]
}

export type IO = {
    inputs: ResourceSet
    outputs: ResourceSet
}

export type Executable = {
    type: string
    data: Record<string, any>
}

export type Specification = {
    id: string
    executable: Executable
    io: IO
    meta: Record<string, any>
}
