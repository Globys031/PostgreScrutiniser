/* tslint:disable */
/* eslint-disable */
/**
 * File Difference API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


import type { Configuration } from '../configuration';
import type { AxiosPromise, AxiosInstance, AxiosRequestConfig } from 'axios';
import globalAxios from 'axios';
// Some imports not used depending on template conditions
// @ts-ignore
import { DUMMY_BASE_URL, assertParamExists, setApiKeyToObject, setBasicAuthToObject, setBearerAuthToObject, setOAuthToObject, setSearchParams, serializeDataIfNeeded, toPathString, createRequestFunction } from '../common';
import type { RequestArgs } from '../base';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, BaseAPI, RequiredError } from '../base';

/**
 * 
 * @export
 * @interface BackupFile
 */
export interface BackupFile {
    /**
     * name of the backup file
     * @type {string}
     * @memberof BackupFile
     */
    'name': string;
    /**
     * timestamp for when the backup file was created
     * @type {string}
     * @memberof BackupFile
     */
    'time': string;
}
/**
 * 
 * @export
 * @interface ErrorMessage
 */
export interface ErrorMessage {
    /**
     * 
     * @type {string}
     * @memberof ErrorMessage
     */
    'error_message': string;
}
/**
 * 
 * @export
 * @interface FileDiffLine
 */
export interface FileDiffLine {
    /**
     * line that is compared between two files
     * @type {string}
     * @memberof FileDiffLine
     */
    'line': string;
    /**
     * specifies whether line has been added, removed or unchanged
     * @type {string}
     * @memberof FileDiffLine
     */
    'type': FileDiffLineTypeEnum;
}

export const FileDiffLineTypeEnum = {
    Equal: 'Equal',
    Insert: 'Insert',
    Delete: 'Delete'
} as const;

export type FileDiffLineTypeEnum = typeof FileDiffLineTypeEnum[keyof typeof FileDiffLineTypeEnum];

/**
 * 
 * @export
 * @interface FileDiffResponse
 */
export interface FileDiffResponse {
    /**
     * 
     * @type {Array<FileDiffLine>}
     * @memberof FileDiffResponse
     */
    'diff': Array<FileDiffLine>;
    /**
     * the name of the file being compared
     * @type {string}
     * @memberof FileDiffResponse
     */
    'filename': string;
    /**
     * timestamp for when the backup file was created
     * @type {string}
     * @memberof FileDiffResponse
     */
    'time': string;
}

/**
 * BackupApi - axios parameter creator
 * @export
 */
export const BackupApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * delete a specific postgresql.auto.conf backup file
         * @param {string} backupName Name of the backup file to delete
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteBackup: async (backupName: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'backupName' is not null or undefined
            assertParamExists('deleteBackup', 'backupName', backupName)
            const localVarPath = `/backup/{backup_name}`
                .replace(`{${"backup_name"}}`, encodeURIComponent(String(backupName)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'DELETE', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication bearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * removes all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteBackups: async (options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/backup`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'DELETE', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication bearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * get all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getBackups: async (options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            const localVarPath = `/backup`;
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication bearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
        /**
         * restores backup file specified by parameter
         * @param {string} backupName file name of the backup to be restored
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        putBackup: async (backupName: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'backupName' is not null or undefined
            assertParamExists('putBackup', 'backupName', backupName)
            const localVarPath = `/backup/{backup_name}`
                .replace(`{${"backup_name"}}`, encodeURIComponent(String(backupName)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'PUT', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication bearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * BackupApi - functional programming interface
 * @export
 */
export const BackupApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = BackupApiAxiosParamCreator(configuration)
    return {
        /**
         * delete a specific postgresql.auto.conf backup file
         * @param {string} backupName Name of the backup file to delete
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async deleteBackup(backupName: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<void>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.deleteBackup(backupName, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * removes all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async deleteBackups(options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<void>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.deleteBackups(options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * get all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getBackups(options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<BackupFile>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getBackups(options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
        /**
         * restores backup file specified by parameter
         * @param {string} backupName file name of the backup to be restored
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async putBackup(backupName: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<void>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.putBackup(backupName, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
    }
};

/**
 * BackupApi - factory interface
 * @export
 */
export const BackupApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = BackupApiFp(configuration)
    return {
        /**
         * delete a specific postgresql.auto.conf backup file
         * @param {string} backupName Name of the backup file to delete
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteBackup(backupName: string, options?: any): AxiosPromise<void> {
            return localVarFp.deleteBackup(backupName, options).then((request) => request(axios, basePath));
        },
        /**
         * removes all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        deleteBackups(options?: any): AxiosPromise<void> {
            return localVarFp.deleteBackups(options).then((request) => request(axios, basePath));
        },
        /**
         * get all postgresql.auto.conf backup files
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getBackups(options?: any): AxiosPromise<BackupFile> {
            return localVarFp.getBackups(options).then((request) => request(axios, basePath));
        },
        /**
         * restores backup file specified by parameter
         * @param {string} backupName file name of the backup to be restored
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        putBackup(backupName: string, options?: any): AxiosPromise<void> {
            return localVarFp.putBackup(backupName, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * BackupApi - object-oriented interface
 * @export
 * @class BackupApi
 * @extends {BaseAPI}
 */
export class BackupApi extends BaseAPI {
    /**
     * delete a specific postgresql.auto.conf backup file
     * @param {string} backupName Name of the backup file to delete
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof BackupApi
     */
    public deleteBackup(backupName: string, options?: AxiosRequestConfig) {
        return BackupApiFp(this.configuration).deleteBackup(backupName, options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * removes all postgresql.auto.conf backup files
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof BackupApi
     */
    public deleteBackups(options?: AxiosRequestConfig) {
        return BackupApiFp(this.configuration).deleteBackups(options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * get all postgresql.auto.conf backup files
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof BackupApi
     */
    public getBackups(options?: AxiosRequestConfig) {
        return BackupApiFp(this.configuration).getBackups(options).then((request) => request(this.axios, this.basePath));
    }

    /**
     * restores backup file specified by parameter
     * @param {string} backupName file name of the backup to be restored
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof BackupApi
     */
    public putBackup(backupName: string, options?: AxiosRequestConfig) {
        return BackupApiFp(this.configuration).putBackup(backupName, options).then((request) => request(this.axios, this.basePath));
    }
}


/**
 * FileDiffApi - axios parameter creator
 * @export
 */
export const FileDiffApiAxiosParamCreator = function (configuration?: Configuration) {
    return {
        /**
         * get difference between current postgresql.auto.conf file and a backup file
         * @param {string} backupName Name of the backup file to diff against
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getFileDiff: async (backupName: string, options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
            // verify required parameter 'backupName' is not null or undefined
            assertParamExists('getFileDiff', 'backupName', backupName)
            const localVarPath = `/file-diff/{backup_name}`
                .replace(`{${"backup_name"}}`, encodeURIComponent(String(backupName)));
            // use dummy base URL string because the URL constructor only accepts absolute URLs.
            const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
            let baseOptions;
            if (configuration) {
                baseOptions = configuration.baseOptions;
            }

            const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options};
            const localVarHeaderParameter = {} as any;
            const localVarQueryParameter = {} as any;

            // authentication bearerAuth required
            // http bearer authentication required
            await setBearerAuthToObject(localVarHeaderParameter, configuration)


    
            setSearchParams(localVarUrlObj, localVarQueryParameter);
            let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
            localVarRequestOptions.headers = {...localVarHeaderParameter, ...headersFromBaseOptions, ...options.headers};

            return {
                url: toPathString(localVarUrlObj),
                options: localVarRequestOptions,
            };
        },
    }
};

/**
 * FileDiffApi - functional programming interface
 * @export
 */
export const FileDiffApiFp = function(configuration?: Configuration) {
    const localVarAxiosParamCreator = FileDiffApiAxiosParamCreator(configuration)
    return {
        /**
         * get difference between current postgresql.auto.conf file and a backup file
         * @param {string} backupName Name of the backup file to diff against
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        async getFileDiff(backupName: string, options?: AxiosRequestConfig): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<FileDiffResponse>> {
            const localVarAxiosArgs = await localVarAxiosParamCreator.getFileDiff(backupName, options);
            return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
        },
    }
};

/**
 * FileDiffApi - factory interface
 * @export
 */
export const FileDiffApiFactory = function (configuration?: Configuration, basePath?: string, axios?: AxiosInstance) {
    const localVarFp = FileDiffApiFp(configuration)
    return {
        /**
         * get difference between current postgresql.auto.conf file and a backup file
         * @param {string} backupName Name of the backup file to diff against
         * @param {*} [options] Override http request option.
         * @throws {RequiredError}
         */
        getFileDiff(backupName: string, options?: any): AxiosPromise<FileDiffResponse> {
            return localVarFp.getFileDiff(backupName, options).then((request) => request(axios, basePath));
        },
    };
};

/**
 * FileDiffApi - object-oriented interface
 * @export
 * @class FileDiffApi
 * @extends {BaseAPI}
 */
export class FileDiffApi extends BaseAPI {
    /**
     * get difference between current postgresql.auto.conf file and a backup file
     * @param {string} backupName Name of the backup file to diff against
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     * @memberof FileDiffApi
     */
    public getFileDiff(backupName: string, options?: AxiosRequestConfig) {
        return FileDiffApiFp(this.configuration).getFileDiff(backupName, options).then((request) => request(this.axios, this.basePath));
    }
}


