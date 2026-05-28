const { createDefaultPreset } = require("ts-jest");

const tsJestTransformCfg = createDefaultPreset().transform;

/** @type {import('ts-jest').JestConfigWithTsJest} */
module.exports = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  testMatch: ['**/*.test.ts'], 
  clearMocks: true,
  coveragePathIgnorePatterns: [
    '/node_modules/',
    'src/index.ts' 
  ],
};