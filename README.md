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

## compute_groundtruth

~~~
./tests/utils/compute_groundtruth  --data_type float --dist_fn l2 --base_file data/sift/sift_learn.fbin --query_file  data/sift/sift_query.fbin --gt_file data/sift/sift_query_learn_gt100 --K 100

~~~