import type { Config } from "jest";
import { defaults } from "jest-config";

const config: Config = {
    preset: "ts-jest",
    testEnvironment: "node",
    setupFilesAfterEnv: ["jest-extended/all"],
    modulePathIgnorePatterns: ["dist"],
    moduleFileExtensions: [...defaults.moduleFileExtensions, "sysl"],
};

export default config;
