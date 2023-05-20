#!/bin/bash
#SBATCH --job-name=cuda
#SBATCH --partition=dgx2
#SBATCH --cpus-per-task=1
#SBATCH --nodes=1
#SBATCH --ntasks-per-node=1
#SBATCH --output=cuda.out
#SBATCH --error=cuda.err
#SBATCH --mail-type=end
#SBATCH --mail-user=jmhuang_2020@sjtu.edu.cn
#SBATCH --gres=1
ulimit -s unlimited
ulimit -l unlimited
module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0
nvcc cuda.cu -o job -lcublas
./job