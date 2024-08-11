/**
 * Kantoku API
 * Create and execute distributed workflows
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 *
 */

import ApiClient from '../ApiClient';

/**
 * The RecordSorter model module.
 * @module model/RecordSorter
 * @version 1.0.0
 */
class RecordSorter {
    /**
     * Constructs a new <code>RecordSorter</code>.
     * @alias module:model/RecordSorter
     * @param key {String} 
     * @param ordering {String} 
     */
    constructor(key, ordering) { 
        
        RecordSorter.initialize(this, key, ordering);
    }

    /**
     * Initializes the fields of this object.
     * This method is used by the constructors of any subclasses, in order to implement multiple inheritance (mix-ins).
     * Only for internal use.
     */
    static initialize(obj, key, ordering) { 
        obj['key'] = key;
        obj['ordering'] = ordering;
    }

    /**
     * Constructs a <code>RecordSorter</code> from a plain JavaScript object, optionally creating a new instance.
     * Copies all relevant properties from <code>data</code> to <code>obj</code> if supplied or a new instance if not.
     * @param {Object} data The plain JavaScript object bearing properties of interest.
     * @param {module:model/RecordSorter} obj Optional instance to populate.
     * @return {module:model/RecordSorter} The populated <code>RecordSorter</code> instance.
     */
    static constructFromObject(data, obj) {
        if (data) {
            obj = obj || new RecordSorter();

            if (data.hasOwnProperty('key')) {
                obj['key'] = ApiClient.convertToType(data['key'], 'String');
            }
            if (data.hasOwnProperty('ordering')) {
                obj['ordering'] = ApiClient.convertToType(data['ordering'], 'String');
            }
        }
        return obj;
    }


}

/**
 * @member {String} key
 */
RecordSorter.prototype['key'] = undefined;

/**
 * @member {String} ordering
 */
RecordSorter.prototype['ordering'] = undefined;






export default RecordSorter;
