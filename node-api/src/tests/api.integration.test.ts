
import { describe, test, expect, beforeAll } from '@jest/globals';
import request from 'supertest';
import jwt from 'jsonwebtoken';
import app from '../app'; 
import { envs } from '../config/env.validation';

describe('Integration Test - API de Node.js', () => {
    
    let validToken: string;

    
    beforeAll(() => {
        
        validToken = jwt.sign({ authorized: true }, envs.jwtSecret.trim(), { algorithm: 'HS256' });
    });

    afterAll(() => {
        jest.resetAllMocks();
    });
    test('Must be rejected if no token is provided (401 Error)', async () => {
        const response = await request(app)
            .post('/api/statistics')
            .send({ data: "algo" });

        expect(response.status).toBe(401);
        expect(response.body.error).toContain('Access denied');
    });

    test('Must be rejected if the token is invalid', async () => {
        const response = await request(app)
            .post('/api/statistics')
            .set('Authorization', 'Bearer token_falso_inventado')
            .send({ data: "algo" });

        expect(response.status).toBe(401);
        expect(response.body.error).toBe('Invalid or expired token.');
    });

    test('Must process the request with a valid token (200 Status)', async () => {
        
        const mockPayload = {
            rotated: [[0,0,1],[0,5,0],[9,0,0]],
            matrix_q: [[0,0,1],[0,1,0],[1,0,0]],
            matrix_r: [[9,0,0],[0,5,0],[0,0,1]]
        };

        const response = await request(app)
            .post('/api/statistics')
            .set('Authorization', `Bearer ${validToken}`) 
            .send(mockPayload);

        
        expect(response.status).toBe(200);
        expect(response.body).toHaveProperty('max_value');
        expect(response.body).toHaveProperty('min_value');
        expect(response.body).toHaveProperty('total_sum');
        expect(response.body).toHaveProperty('average');
    });

    test('Debe manejar matrices completamente vacías y retornar promedio 0 (Cobertura división por cero)', async () => {
        const emptyPayload = {
            rotated: [],
            matrix_q: [],
            matrix_r: []
        };

        const response = await request(app)
            .post('/api/statistics')
            .set('Authorization', `Bearer ${validToken}`)
            .send(emptyPayload);

        expect(response.status).toBe(200);
        
        expect(response.body.average).toBe(0);
        expect(response.body.total_elements).toBe(undefined);
    });
    test('Debe retornar 400 si faltan matrices en el body (Cobertura error 400)', async () => {
        const response = await request(app)
            .post('/api/statistics')
            .set('Authorization', `Bearer ${validToken}`)
            .send({ rotated: [[1,2],[3,4]] }); 

        expect(response.status).toBe(400);
        expect(response.body.error).toBe('Missing matrices in the request body');
    });

    test('Debe evaluar matrices vacías o no cuadradas como NO diagonales (Cobertura isDiagonalMatrix)', async () => {
        const payloadConTrampa = {
            rotated: [], 
            matrix_q: [[1, 2, 3], [4, 5, 6]], 
            matrix_r: [[1, 0], [0, 1]] 
        };

        const response = await request(app)
            .post('/api/statistics')
            .set('Authorization', `Bearer ${validToken}`)
            .send(payloadConTrampa);

        expect(response.status).toBe(200);
        
        expect(response.body.diagonal_matrix_verification.rotated).toBe(false);
        expect(response.body.diagonal_matrix_verification.matrix_q).toBe(false);
    });
});