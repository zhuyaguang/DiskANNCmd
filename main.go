package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	binPath  string
	QfvecPath string
	QfbinPath string
	LfvecPath string
	LfbinPath string

	dataType  string
	distFn    string
	baseFile  string
	queryFile string
	gtFile    string
	K         string

	dataPath        string
	indexPathPrefix string

	resultK         string
	L               string
	resultPath      string
	numNodesToCache string
)

func init() {
	// vec_to_bin
	binPath = getEnvOrDefault("BIN_PATH", "/home/zjlab/zyg/bin/")
	QfvecPath = getEnvOrDefault("QFVEC_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fvecs")
	QfbinPath = getEnvOrDefault("QFBIN_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fbin")
	LfvecPath = getEnvOrDefault("LFVEC_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fvecs")
	LfbinPath = getEnvOrDefault("LFBIN_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")

	// compute_groundtruth
	dataType = getEnvOrDefault("DATA_TYPE", "float")
	distFn = getEnvOrDefault("DIST_FN", "l2")
	baseFile = getEnvOrDefault("BASE_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")
	queryFile = getEnvOrDefault("QUERY_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fbin")
	gtFile = getEnvOrDefault("GT_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query_learn_gt100")
	K = getEnvOrDefault("K", "100")

	// BuildDiskIndex
	dataPath = getEnvOrDefault("DATA_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")
	indexPathPrefix = getEnvOrDefault("INDEX_PATH_PREFIX", "/home/zjlab/zyg/DiskANN/build/data/sift/disk_index_sift_learn_R32_L50_A1.2")

	// SearchDiskIndex
	resultK = getEnvOrDefault("RESULT_K", "10")
	L = getEnvOrDefault("L", "10")
	resultPath = getEnvOrDefault("RESULT_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/res")
	numNodesToCache = getEnvOrDefault("NUM_NODES_TO_CACHE", "10000")

	// 初始化原始数据和构建索引
	LearnVecToBin()
	LearnBiludIndex()
}

func main() {
	flag.Parse()

	router := gin.Default()
	router.POST("/VecToBin", postVecToBin)
	router.POST("/SearchDiskIndex", postSearchDiskIndex)

	router.Run(":18180")
}

func LearnVecToBin()  {
	err := FvecToBin(binPath, LfvecPath, LfbinPath)
	if err != nil {
		return
	}
}

func LearnBiludIndex()  {
	err := BuildDiskIndex(binPath, dataType, distFn, dataPath, indexPathPrefix)
	if err != nil {
		return
	}
}

func postVecToBin(c *gin.Context) {

	err := FvecToBin(binPath, QfvecPath, QfbinPath)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusCreated, "VecToBin successful")
}


func postSearchDiskIndex(c *gin.Context) {
	start := time.Now()
	// Code to measure
	// 1.text to vec

	// 2.vec to bin

	err := FvecToBin(binPath, QfvecPath, QfbinPath)
	if err != nil {
		return
	}
	duration := time.Since(start)
	fmt.Println(duration)

	// 3.postComputeGroundTruth 可省略

	// 3.SearchDiskIndex

	err ,_ = SearchDiskIndex(binPath, dataType, distFn, indexPathPrefix, queryFile, gtFile, resultK, L, resultPath, numNodesToCache)
	if err != nil {
		return
	}
	duration = time.Since(start)
	fmt.Println(duration)
	c.IndentedJSON(http.StatusCreated, "SearchDiskIndex successful")
}

func getEnvOrDefault(env string, defaultValue string) string {
	v := os.Getenv(env)
	if v == "" {
		v = defaultValue
	}
	return v
}

func FvecToBin(bin, fvecPath, fbinPath string) error {
	prg := bin + "fvecs_to_bin"
	cmd := exec.Command(prg, fvecPath, fbinPath)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print("FvecToBin:", string(stdout))
	return nil
}

// ComputeGroundTruth  ./tests/utils/compute_groundtruth  --data_type float --dist_fn l2 --base_file data/sift/sift_learn.fbin --query_file  data/sift/sift_query.fbin --gt_file data/sift/sift_query_learn_gt100 --K 100
func ComputeGroundTruth(bin, dataType, distFn, baseFile, queryFile, gtFile, K string) error {
	prg := bin + "compute_groundtruth"
	cmdString := "--data_type " + dataType + " --dist_fn " + distFn + " --base_file " + baseFile + " --query_file " + queryFile + " --gt_file " + gtFile + " --K" + K
	fmt.Println(cmdString)
	cmd := exec.Command("sh", "-c", prg +" "+ cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print("ComputeGroundTruth:", string(stdout))
	return nil
}

// BuildDiskIndex  ./tests/build_disk_index --data_type float --dist_fn l2 --data_path data/sift/sift_learn.fbin --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 -R 32 -L50 -B 0.003 -M 1
func BuildDiskIndex(bin, dataType, distFn, dataPath, indexPathPrefix string) error {
	prg := bin + "build_disk_index"
	subCmd := " -R 32 -L50 -B 0.003 -M 1"
	cmdString := fmt.Sprintf("--data_type " + dataType + " --dist_fn " + distFn + " --data_path " + dataPath + " --index_path_prefix " + indexPathPrefix + subCmd)
	cmd := exec.Command("sh", "-c", prg +" "+ cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print("BuildDiskIndex:", string(stdout))

	return nil
}

// SearchDiskIndex  ./tests/search_disk_index  --data_type float --dist_fn l2 --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 --query_file data/sift/sift_query.fbin  --gt_file data/sift/sift_query_learn_gt100 -K 10 -L 10 20 30 40 50 100 --result_path data/sift/res --num_nodes_to_cache 10000
func SearchDiskIndex(bin, dataType, distFn, indexPathPrefix, queryFile, gtFile, K, L, resultPath, numNodesToCache string) (error,string) {
	prg := bin + "search_disk_index"
	cmdString := fmt.Sprintf("--data_type " + dataType + " --dist_fn " + distFn + " --index_path_prefix " + indexPathPrefix + " --query_file " + queryFile + " --gt_file " + gtFile + " -K " + K + " -L " + L + " --result_path " + resultPath + " --num_nodes_to_cache " + numNodesToCache)
	cmd := exec.Command("sh", "-c", prg +" "+ cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err,""
	}

	fmt.Print("SearchDiskIndex:", string(stdout))

	return nil,string(stdout)
}
