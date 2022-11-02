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
                                   --data_type float --dist_fn l2 --base_file /home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin --query_file /home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fbin --gt_file /home/zjlab/zyg/DiskANN/build/data/sift/sift_query_learn_gt100 --K 100
./tests/build_disk_index --data_type float --dist_fn l2 --data_path data/sift/sift_learn.fbin --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 -R 32 -L50 -B 0.003 -
M 1
 ./tests/search_disk_index  --data_type float --dist_fn l2 --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 --query_file data/sift/sift_query.fbin  --gt_file data/sift/sift_query_learn_gt100 -K 10 -L 10 20 30 40 50 100 --result_path data/sift/res --num_nodes_to_cache 10000
 
~~~