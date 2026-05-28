import React, { useState } from 'react';
import { processMatrix } from './services/api.service';
import type { MatrixResponse } from './types/matrix';

type MatrixGridProps = {
  title: string;
  matrix: number[][];
};

function MatrixGrid({ title, matrix }: MatrixGridProps) {
  return (
    <div className="bg-white p-4 rounded-lg shadow border border-gray-100">
      <h3 className="font-bold text-gray-700 mb-2">{title}</h3>

      <div className="font-mono text-sm bg-gray-50 p-2 rounded">
        {matrix.map((row, i) => (
          <div key={i} className="flex space-x-2">
            {row.map((val, j) => (
              <span key={j} className="w-12 text-right">
                {Number(val).toFixed(2)}
              </span>
            ))}
          </div>
        ))}
      </div>
    </div>
  );
}

function App() {
  const [inputStr, setInputStr] = useState(
    '[\n  [1, 2, 3],\n  [4, 5, 6],\n  [7, 8, 9]\n]'
  );

  const [data, setData] = useState<MatrixResponse | null>(null);

  const [loading, setLoading] = useState(false);

  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (
    e: React.SubmitEvent<HTMLFormElement>
  ) => {
    e.preventDefault();

    setLoading(true);
    setError(null);
    setData(null);

    try {
      const parsedMatrix: number[][] = JSON.parse(inputStr);

      if (
        !Array.isArray(parsedMatrix) ||
        !Array.isArray(parsedMatrix[0])
      ) {
        throw new Error(
          'Formato inválido. Debe ser un arreglo 2D de números.'
        );
      }

      const result = await processMatrix(parsedMatrix);

      if (!result) {
        throw new Error('No se obtuvo respuesta del servidor');
      }

      setData(result);

    } catch (error: unknown) {
      if (error instanceof Error) {
        setError(error.message);
      } else {
        setError('Error procesando el JSON');
      }

    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-5xl mx-auto space-y-6">

        <header className="text-center">
          <h1 className="text-3xl font-bold text-slate-800">
            Procesador QR Interseguro
          </h1>

          <p className="text-slate-500">
            Go (Rotación & QR) + Node.js (Estadísticas)
          </p>
        </header>

        <section className="bg-white p-6 rounded-xl shadow-md">
          <form
            onSubmit={handleSubmit}
            className="flex flex-col space-y-4"
          >
            <label className="font-semibold text-slate-700">
              Ingresa la matriz (JSON):
            </label>

            <textarea
              className="w-full h-32 p-3 font-mono text-sm border border-slate-300 rounded focus:ring-2 focus:ring-blue-500 outline-none"
              value={inputStr}
              onChange={(e) => setInputStr(e.target.value)}
            />

            <button
              type="submit"
              disabled={loading}
              className="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-300 text-white font-bold py-2 px-4 rounded transition-all"
            >
              {loading
                ? 'Procesando en el backend...'
                : 'Calcular Matrices'}
            </button>
          </form>

          {error && (
            <div className="mt-4 p-3 bg-red-100 text-red-700 rounded font-medium">
              {error}
            </div>
          )}
        </section>

        {data && (
          <div className="space-y-6 animate-fade-in">

            <div>
              <h2 className="text-xl font-bold text-slate-800 mb-3 border-b pb-2">
                Resultados en Go
              </h2>

              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <MatrixGrid
                  title="Rotada 90°"
                  matrix={data.go_results.rotated}
                />

                <MatrixGrid
                  title="Matriz Q"
                  matrix={data.go_results.matrix_q}
                />

                <MatrixGrid
                  title="Matriz R"
                  matrix={data.go_results.matrix_r}
                />
              </div>
            </div>
           
            <div>
              <h2 className="text-xl font-bold text-slate-800 mb-3 border-b pb-2">
                Estadísticas en Node.js
              </h2>

              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">

                <div className="bg-slate-800 p-4 rounded-lg shadow">
                  <p className="text-slate-400 text-sm">
                    Valor Máximo
                  </p>

                  <p className="text-2xl font-bold text-white">
                    {data.node_stats.max_value}
                  </p>
                </div>

                <div className="bg-slate-800 p-4 rounded-lg shadow">
                  <p className="text-slate-400 text-sm">
                    Valor Mínimo
                  </p>

                  <p className="text-2xl font-bold text-white">
                    {data.node_stats.min_value}
                  </p>
                </div>

                <div className="bg-slate-800 p-4 rounded-lg shadow">
                  <p className="text-slate-400 text-sm">
                    Suma Total
                  </p>

                  <p className="text-2xl font-bold text-white">
                    {data.node_stats.total_sum}
                  </p>
                </div>

                <div className="bg-slate-800 p-4 rounded-lg shadow">
                  <p className="text-slate-400 text-sm">
                    Promedio
                  </p>

                  <p className="text-2xl font-bold text-white">
                    {data.node_stats.average}
                  </p>
                </div>
                <div>
            </div>
              </div>
            </div>
             <div>
            <div className="mt-8">
              <h2 className="text-xl font-bold text-slate-800 mb-3 border-b pb-2">
                Verificación de Matrices Diagonales
              </h2>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-center">

                <div className={`p-4 rounded-lg shadow border ${data.node_stats.diagonal_matrix_verification.rotated ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
                  <p className="text-slate-600 text-sm font-semibold mb-1">Rotada 90°</p>
                  <p className={`text-xl font-bold ${data.node_stats.diagonal_matrix_verification.rotated ? 'text-green-700' : 'text-red-700'}`}>
                    {data.node_stats.diagonal_matrix_verification.rotated ? '✅ Sí' : '❌ No'}
                  </p>
                </div>

                <div className={`p-4 rounded-lg shadow border ${data.node_stats.diagonal_matrix_verification.matrix_q ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
                  <p className="text-slate-600 text-sm font-semibold mb-1">Matriz Q</p>
                  <p className={`text-xl font-bold ${data.node_stats.diagonal_matrix_verification.matrix_q ? 'text-green-700' : 'text-red-700'}`}>
                    {data.node_stats.diagonal_matrix_verification.matrix_q ? '✅ Sí' : '❌ No'}
                  </p>
                </div>

                <div className={`p-4 rounded-lg shadow border ${data.node_stats.diagonal_matrix_verification.matrix_r ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'}`}>
                  <p className="text-slate-600 text-sm font-semibold mb-1">Matriz R</p>
                  <p className={`text-xl font-bold ${data.node_stats.diagonal_matrix_verification.matrix_r ? 'text-green-700' : 'text-red-700'}`}>
                    {data.node_stats.diagonal_matrix_verification.matrix_r ? '✅ Sí' : '❌ No'}
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

      </div>
    </div>
  );
}

export default App;