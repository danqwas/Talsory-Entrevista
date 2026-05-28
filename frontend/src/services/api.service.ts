import type { MatrixResponse } from '../types/matrix';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api';

type LoginResponse = {
  token: string;
};

type ApiError = {
  error?: string;
};

async function parseResponse<T>(response: Response): Promise<T> {
  const data = (await response.json()) as T;

  if (!response.ok) {
    const error = data as ApiError;
    throw new Error(error.error || 'Error en la solicitud');
  }

  return data;
}

export async function processMatrix(
  matrix: number[][]
): Promise<MatrixResponse | undefined> {
  try {
    const loginResponse = await fetch(`${API_URL}/auth/login`, {
      method: 'POST',
    });

    const { token } =
      await parseResponse<LoginResponse>(loginResponse);

    const processResponse = await fetch(
      `${API_URL}/matrix/process`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ matrix }),
      }
    );

    return await parseResponse<MatrixResponse>(processResponse);

  } catch (error: unknown) {
    if (error instanceof Error) {
      throw new Error(error.message, { cause: error });
    }
  }
}