# DiskANNCmd

## 向量转换成二进制

~~~
mkdir -p DiskANN/build/data && cd DiskANN/build/data
wget ftp://ftp.irisa.fr/local/texmex/corpus/sift.tar.gz
tar -xf sift.tar.gz
cd ..
./tests/utils/fvecs_to_bin data/sift/sift_learn.fvecs data/sift/sift_learn.fbin
./tests/utils/fvecs_to_bin data/sift/sift_query.fvecs data/sift/sift_query.fbin
~~~


## 

export BASE_PATH="/home/zjlab/zyg" 

export VECINIT_PATH="/home/zjlab/zyg"

[//]: # (export LFVEC_PATH="/home/zjlab/zyg/vec-init/vectors/init/name.vec")

[//]: # (export LFBIN_PATH="/home/zjlab/zyg/vec-init/vectors/init/name.bin")

