#include <cuda_runtime.h>
#include <iostream>
#include <vector>
#include <fstream>

using namespace std;

__global__ void multi_gpu(int *d_matrixA, int *d_matrixB, int *d_matrixC, int x, int y)
{
    int ix = threadIdx.x + blockDim.x * blockIdx.x;
    int iy = threadIdx.y + blockDim.y * blockIdx.y;
    unsigned int idx = iy * x + ix;
    if (ix < x && iy < y)
    {
        int k;
        int sum = 0;
        for (k = 0; k < x; k++)
        {
            sum += d_matrixA[iy * x + k] * d_matrixB[k * x + ix];
        }
        d_matrixC[idx] = sum;
    }
}

vector<vector<int>> matrix_add(vector<vector<int>> &a, vector<vector<int>> &b)
{
    const int m = a.size(), n = a[0].size();
    int *matrixA = (int *)malloc(sizeof(int) * m * n);
    int *matrixB = (int *)malloc(sizeof(int) * m * n);
    int *matrixC = (int *)malloc(sizeof(int) * m * n);

    for (int i = 0; i < m; i++)
    {
        for (int j = 0; j < n; j++)
        {
            matrixA[i * n + j] = a[i][j];
            matrixB[i * n + j] = b[i][j];
        }
    }

    int *d_matrixA, *d_matrixB, *d_matrixC;
    cudaMalloc((void **)&d_matrixA, sizeof(int) * n * m);
    cudaMalloc((void **)&d_matrixB, sizeof(int) * n * m);
    cudaMalloc((void **)&d_matrixC, sizeof(int) * n * m);
    cudaMemcpy(d_matrixA, matrixA, sizeof(int) * n * m, cudaMemcpyHostToDevice);
    cudaMemcpy(d_matrixB, matrixB, sizeof(int) * n * m, cudaMemcpyHostToDevice);
    int x = n, y = m;
    int dimx = 32;
    int dimy = 32;
    dim3 block(dimx, dimy);
    dim3 grid(x / block.x + 1, y / block.y + 1);
    multi_gpu<<<grid, block>>>(d_matrixA, d_matrixB, d_matrixC, x, y);
    cudaMemcpy(matrixC, d_matrixC, sizeof(int) * n * m, cudaMemcpyDeviceToHost);
    vector<int> temp(n, 0);
    vector<vector<int>> c(m, temp);
    for (int i = 0; i < m; i++)
    {
        for (int j = 0; j < n; j++)
        {
            c[i][j] = matrixC[i * n + j];
        }
    }
    free(matrixA);
    free(matrixB);
    free(matrixC);
    cudaFree(d_matrixA);
    cudaFree(d_matrixB);
    cudaFree(d_matrixC);
    return c;
}

int main(int argc, char *argv[])
{
    vector<vector<int>> a{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, b{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}};
    vector<vector<int>> c = matrix_add(a, b);
    for (int i = 0; i < c.size(); i++)
    {
        for (int j = 0; j < c[0].size(); j++)
        {
            std::cout << c[i][j] << " ";
        }
        std::cout << "\n";
    }

    return 0;
}
