import { describe, test, expect, beforeAll } from '@jest/globals';
jest.mock('dotenv/config', () => jest.fn());

describe('Validación de Variables de Entorno (Joi)', () => {
    const originalEnv = process.env;

    beforeEach(() => {
        jest.resetModules(); 
        process.env = { ...originalEnv }; 
    });

    afterEach(() => {
        process.env = originalEnv; 
    });

    test('Debe lanzar un error si falta el JWT_SECRET', () => {
        
        delete process.env.JWT_SECRET;
        
        expect(() => {
            require('../config/env.validation');
        }).toThrow(/Config validation error/);
    });
});1