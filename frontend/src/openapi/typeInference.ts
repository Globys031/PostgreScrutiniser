import type { ResourceConfig } from "./api/resource-config";

// Transforms snake case words to pascal case
export type ToPascalCase<S extends string> = S extends `${infer T}_${infer U}`
  ? `${Capitalize<T>}${ToPascalCase<U>}`
  : Capitalize<S>;

/* 
Transforms `ResourceConfig` key names to pascal case
Axios response is an object with PascalCase keys. To avoid
creating a custom ResourceConfig interface, this type transforms `ResourceConfig`.
*/
export type ResourceConfigPascalCase = {
  [K in keyof ResourceConfig as ToPascalCase<K & string>]: ResourceConfig[K];
};
