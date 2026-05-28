import express, {Request, Response} from 'express';
import { envs } from './config/env.validation';
import { verifyJWTMiddleware } from './middlewares/auth.middleware';
const app = express();
const PORT = envs.port;

app.use(express.json());


interface MatrixRequestBody {
    rotated: number[][];
    matrix_q: number[][];
    matrix_r: number[][];
}


function isDiagonalMatrix(matrix: number[][]): boolean {
    if (!matrix || matrix.length === 0 || matrix.length !== matrix[0].length) {
        return false; 
    }
    
    const rows = matrix.length;
    const EPSILON = 1e-9; 

    for (let i = 0; i < rows; i++) {
        for (let j = 0; j < rows; j++) {
            if (i !== j && Math.abs(matrix[i][j]) > EPSILON) {
                return false;
            }
        }
    }
    return true;
}


app.post('/api/statistics', verifyJWTMiddleware, (req: Request<{}, {}, MatrixRequestBody>, res: Response) => {
    const { rotated, matrix_q, matrix_r } = req.body;

    if (!rotated || !matrix_q || !matrix_r) {
          return res.status(400).json({ error: 'Missing matrices in the request body' });
      }

    
    const matrices: Record<string, number[][]> = { rotated, matrix_q, matrix_r };
    
    let maxVal = -Infinity;
    let minVal = Infinity;
    let totalSum = 0;
    let totalElements = 0;
    
    
    const diagonalChecks: Record<string, boolean> = {};

    for (const [key, matrix] of Object.entries(matrices)) {
        diagonalChecks[key] = isDiagonalMatrix(matrix);

        for (let i = 0; i < matrix.length; i++) {
            for (let j = 0; j < matrix[i].length; j++) {
                const val = matrix[i][j];
                
                if (val > maxVal) maxVal = val;
                if (val < minVal) minVal = val;
                
                totalSum += val;
                totalElements++;
            }
        }
    }

    const average = totalElements > 0 ? totalSum / totalElements : 0;

    return res.json({
        max_value: maxVal,
        min_value: minVal,
        total_sum: totalSum,
        average: Number(average.toFixed(4)),
        diagonal_matrix_verification: diagonalChecks
    });
});


export default app;