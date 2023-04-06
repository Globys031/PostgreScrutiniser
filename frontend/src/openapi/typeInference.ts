import type { ResourceConfig } from "./api/resource-config";

// Transforms snake case words to pascal case
type ToPascalCase<S extends string> = S extends `${infer T}_${infer U}`
  ? `${Capitalize<T>}${ToPascalCase<U>}`
  : Capitalize<S>;

// Transforms `<T>` key names to pascal case. Axios response is an object with PascalCase keys.
// To avoid creating a custom ResourceConfig interface, this type transforms `ResourceConfig`.
type PascalCase<T> = {
  [K in keyof T as ToPascalCase<K & string>]: T[K];
};

export type ResourceConfigPascalCase = PascalCase<ResourceConfig>;
