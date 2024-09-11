# KantokuApi.DefaultApi

All URIs are relative to *https://kantoku.api.server:3000*

Method | HTTP request | Description
------------- | ------------- | -------------
[**resourcesAllocatePost**](DefaultApi.md#resourcesAllocatePost) | **POST** /resources/allocate | Allocates N resources
[**resourcesDeallocatePost**](DefaultApi.md#resourcesDeallocatePost) | **POST** /resources/deallocate | Deallocate resources
[**resourcesInitializePost**](DefaultApi.md#resourcesInitializePost) | **POST** /resources/initialize | Initialize resources
[**resourcesLoadPost**](DefaultApi.md#resourcesLoadPost) | **POST** /resources/load | Load resources
[**tasksCountPost**](DefaultApi.md#tasksCountPost) | **POST** /tasks/count | Count records using a filter
[**tasksFilterPost**](DefaultApi.md#tasksFilterPost) | **POST** /tasks/filter | Load records using a filter
[**tasksLoadPost**](DefaultApi.md#tasksLoadPost) | **POST** /tasks/load | Load a set of tasks
[**tasksSpawnFromSpecPost**](DefaultApi.md#tasksSpawnFromSpecPost) | **POST** /tasks/spawn_from_spec | Spawn a new task from specification
[**tasksSpawnPost**](DefaultApi.md#tasksSpawnPost) | **POST** /tasks/spawn | Spawn a new task
[**tasksSpecificationsCreatePost**](DefaultApi.md#tasksSpecificationsCreatePost) | **POST** /tasks/specifications/create | Create a specification
[**tasksSpecificationsGetAllPost**](DefaultApi.md#tasksSpecificationsGetAllPost) | **POST** /tasks/specifications/get_all | Get all specifications
[**tasksSpecificationsGetPost**](DefaultApi.md#tasksSpecificationsGetPost) | **POST** /tasks/specifications/get | Get specifications by id
[**tasksSpecificationsRemovePost**](DefaultApi.md#tasksSpecificationsRemovePost) | **POST** /tasks/specifications/remove | Remove a specification
[**tasksSpecificationsTypesCreatePost**](DefaultApi.md#tasksSpecificationsTypesCreatePost) | **POST** /tasks/specifications/types/create | Create a type
[**tasksSpecificationsTypesGetAllPost**](DefaultApi.md#tasksSpecificationsTypesGetAllPost) | **POST** /tasks/specifications/types/get_all | Get all types
[**tasksSpecificationsTypesGetPost**](DefaultApi.md#tasksSpecificationsTypesGetPost) | **POST** /tasks/specifications/types/get | Get a type by id
[**tasksSpecificationsTypesRemovePost**](DefaultApi.md#tasksSpecificationsTypesRemovePost) | **POST** /tasks/specifications/types/remove | Remove a type
[**tasksUpdatePost**](DefaultApi.md#tasksUpdatePost) | **POST** /tasks/update | Update a record



## resourcesAllocatePost

> [String] resourcesAllocatePost(amount)

Allocates N resources

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let amount = 56; // Number | 
apiInstance.resourcesAllocatePost(amount, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **amount** | **Number**|  | 

### Return type

**[String]**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## resourcesDeallocatePost

> Object resourcesDeallocatePost(requestBody)

Deallocate resources

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let requestBody = ["null"]; // [String] | A list of resource_db identifiers
apiInstance.resourcesDeallocatePost(requestBody, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestBody** | [**[String]**](String.md)| A list of resource identifiers | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## resourcesInitializePost

> Object resourcesInitializePost(resourceInitializer)

Initialize resources

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let resourceInitializer = [new KantokuApi.ResourceInitializer()]; // [ResourceInitializer] | A dictionary (ResourceID -> Value)
apiInstance.resourcesInitializePost(resourceInitializer, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **resourceInitializer** | [**[ResourceInitializer]**](ResourceInitializer.md)| A dictionary (ResourceID -&gt; Value) | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## resourcesLoadPost

> [Resource] resourcesLoadPost(requestBody)

Load resources

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let requestBody = ["null"]; // [String] | A list of resource_db identifiers
apiInstance.resourcesLoadPost(requestBody, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestBody** | [**[String]**](String.md)| A list of resource identifiers | 

### Return type

[**[Resource]**](Resource.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksCountPost

> Number tasksCountPost(tasksFilterPostRequest)

Count records using a filter

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksFilterPostRequest = new KantokuApi.TasksFilterPostRequest(); // TasksFilterPostRequest | A query
apiInstance.tasksCountPost(tasksFilterPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksFilterPostRequest** | [**TasksFilterPostRequest**](TasksFilterPostRequest.md)| A query | 

### Return type

**Number**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksFilterPost

> [Task] tasksFilterPost(tasksFilterPostRequest)

Load records using a filter

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksFilterPostRequest = new KantokuApi.TasksFilterPostRequest(); // TasksFilterPostRequest | A query
apiInstance.tasksFilterPost(tasksFilterPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksFilterPostRequest** | [**TasksFilterPostRequest**](TasksFilterPostRequest.md)| A query | 

### Return type

[**[Task]**](Task.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksLoadPost

> [Task] tasksLoadPost(requestBody)

Load a set of tasks

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let requestBody = ["null"]; // [String] | A list of task identifiers
apiInstance.tasksLoadPost(requestBody, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **requestBody** | [**[String]**](String.md)| A list of task identifiers | 

### Return type

[**[Task]**](Task.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpawnFromSpecPost

> TaskSpawnResponse tasksSpawnFromSpecPost(specificationBasedTaskParameters)

Spawn a new task from specification

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let specificationBasedTaskParameters = new KantokuApi.SpecificationBasedTaskParameters(); // SpecificationBasedTaskParameters | The specification of a task to be spawned
apiInstance.tasksSpawnFromSpecPost(specificationBasedTaskParameters, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **specificationBasedTaskParameters** | [**SpecificationBasedTaskParameters**](SpecificationBasedTaskParameters.md)| The specification of a task to be spawned | 

### Return type

[**TaskSpawnResponse**](TaskSpawnResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpawnPost

> TaskSpawnResponse tasksSpawnPost(taskParameters)

Spawn a new task

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let taskParameters = new KantokuApi.TaskParameters(); // TaskParameters | The specification of a task to be spawned
apiInstance.tasksSpawnPost(taskParameters, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **taskParameters** | [**TaskParameters**](TaskParameters.md)| The specification of a task to be spawned | 

### Return type

[**TaskSpawnResponse**](TaskSpawnResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsCreatePost

> tasksSpecificationsCreatePost(specification)

Create a specification

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let specification = new KantokuApi.Specification(); // Specification | 
apiInstance.tasksSpecificationsCreatePost(specification, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **specification** | [**Specification**](Specification.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsGetAllPost

> [Specification] tasksSpecificationsGetAllPost()

Get all specifications

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
apiInstance.tasksSpecificationsGetAllPost((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**[Specification]**](Specification.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## tasksSpecificationsGetPost

> Specification tasksSpecificationsGetPost(tasksSpecificationsGetPostRequest)

Get specifications by id

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksSpecificationsGetPostRequest = new KantokuApi.TasksSpecificationsGetPostRequest(); // TasksSpecificationsGetPostRequest | 
apiInstance.tasksSpecificationsGetPost(tasksSpecificationsGetPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksSpecificationsGetPostRequest** | [**TasksSpecificationsGetPostRequest**](TasksSpecificationsGetPostRequest.md)|  | 

### Return type

[**Specification**](Specification.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsRemovePost

> tasksSpecificationsRemovePost(tasksSpecificationsGetPostRequest)

Remove a specification

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksSpecificationsGetPostRequest = new KantokuApi.TasksSpecificationsGetPostRequest(); // TasksSpecificationsGetPostRequest | 
apiInstance.tasksSpecificationsRemovePost(tasksSpecificationsGetPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksSpecificationsGetPostRequest** | [**TasksSpecificationsGetPostRequest**](TasksSpecificationsGetPostRequest.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsTypesCreatePost

> tasksSpecificationsTypesCreatePost(typeWithID)

Create a type

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let typeWithID = new KantokuApi.TypeWithID(); // TypeWithID | 
apiInstance.tasksSpecificationsTypesCreatePost(typeWithID, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **typeWithID** | [**TypeWithID**](TypeWithID.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsTypesGetAllPost

> [TypeWithID] tasksSpecificationsTypesGetAllPost()

Get all types

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
apiInstance.tasksSpecificationsTypesGetAllPost((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**[TypeWithID]**](TypeWithID.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## tasksSpecificationsTypesGetPost

> TypeWithID tasksSpecificationsTypesGetPost(tasksSpecificationsGetPostRequest)

Get a type by id

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksSpecificationsGetPostRequest = new KantokuApi.TasksSpecificationsGetPostRequest(); // TasksSpecificationsGetPostRequest | 
apiInstance.tasksSpecificationsTypesGetPost(tasksSpecificationsGetPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksSpecificationsGetPostRequest** | [**TasksSpecificationsGetPostRequest**](TasksSpecificationsGetPostRequest.md)|  | 

### Return type

[**TypeWithID**](TypeWithID.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksSpecificationsTypesRemovePost

> tasksSpecificationsTypesRemovePost(tasksSpecificationsGetPostRequest)

Remove a type

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksSpecificationsGetPostRequest = new KantokuApi.TasksSpecificationsGetPostRequest(); // TasksSpecificationsGetPostRequest | 
apiInstance.tasksSpecificationsTypesRemovePost(tasksSpecificationsGetPostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksSpecificationsGetPostRequest** | [**TasksSpecificationsGetPostRequest**](TasksSpecificationsGetPostRequest.md)|  | 

### Return type

null (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## tasksUpdatePost

> Object tasksUpdatePost(tasksUpdatePostRequest)

Update a record

### Example

```javascript
import KantokuApi from 'kantoku_api';

let apiInstance = new KantokuApi.DefaultApi();
let tasksUpdatePostRequest = new KantokuApi.TasksUpdatePostRequest(); // TasksUpdatePostRequest | 
apiInstance.tasksUpdatePost(tasksUpdatePostRequest, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tasksUpdatePostRequest** | [**TasksUpdatePostRequest**](TasksUpdatePostRequest.md)|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

