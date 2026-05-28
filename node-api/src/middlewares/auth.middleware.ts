import { Request, Response, NextFunction } from 'express';
import jwt from 'jsonwebtoken';
import { envs } from '../config/env.validation';

export function verifyJWTMiddleware(req: Request, res: Response, next: NextFunction) {
    const authHeader = req.headers['authorization'];
    
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({ error: 'Access denied. Bearer token not provided.' });
    }

    const token = authHeader.split(' ')[1];

    try {
        const secretLimpio = envs.jwtSecret.trim();
        const decoded = jwt.verify(token, secretLimpio, { algorithms: ['HS256'] });
        
        (req as any).user = decoded; 
        next();
    } catch (error) {
        return res.status(401).json({ error: 'Invalid or expired token.' });
    }
}