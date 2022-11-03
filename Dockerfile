FROM ubuntu:20.04

RUN mkdir -p /usr/diskann/bin
RUN mkdir -p /usr/diskann/file
COPY bin /usr/diskann/bin
COPY file /usr/diskann/file
RUN chmod +x /usr/diskann/bin/*
# 二进制文件地址
ENV BIN_PATH=/usr/diskann/bin/

# 原始数据向量地址
ENV LFVEC_PATH=/usr/diskann/file/learn.fvecs
# 原始数据二进制地址
ENV LFBIN_PATH=/usr/diskann/file/learn.fbin
ENV DATA_PATH=/usr/diskann/file/learn.fbin
ENV INDEX_PATH_PREFIX=/usr/diskann/file/disk_index_sift_learn_R32_L50_A1.2
ENV DIST_FN=l2
ENV DATA_TYPE=float

# 查询向量地址
ENV QFVEC_PATH=/usr/diskann/file/query.fvecs
# 查询二进制地址
ENV QFBIN_PATH=/usr/diskann/file/query.fbin

EXPOSE 18180
CMD /usr/diskann/bin/main