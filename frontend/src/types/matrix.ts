export interface MatrixResponse {
  go_results: {
    rotated: number[][];
    matrix_q: number[][];
    matrix_r: number[][];
  };
  node_stats: {
    max_value: number;
    min_value: number;
    total_sum: number;
    average: number;
    diagonal_matrix_verification: {
      rotated: boolean;
      matrix_q: boolean;
      matrix_r: boolean;
    };
  };
}
