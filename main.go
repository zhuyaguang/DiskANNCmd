package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	binPath     string
	BasePath    string
	VecInitPath string //vec-init 项目地址
	LfvecPath   string
	LfbinPath   string

	dataType  string
	distFn    string
	baseFile  string
	queryFile string
	gtFile    string
	K         string

	indexPathPrefix string

	resultK         string
	L               string
	resultPath      string
	numNodesToCache string

	healthState int
)

type VecToBin struct {
	Fvec string `json:"fvec"`
	Fbin string `json:"fbin"`
}

func init() {
	// vec_to_bin
	binPath = getEnvOrDefault("BIN_PATH", "/home/zjlab/zyg/bin/")
	BasePath = getEnvOrDefault("BASE_PATH", "/home/zjlab/zyg/")
	VecInitPath = getEnvOrDefault("VECINIT_PATH", "/home/zjlab/zyg/")
	LfvecPath = filepath.Join(VecInitPath, "vectors/init/name.vec")
	LfbinPath = filepath.Join(VecInitPath, "vectors/init/name.bin")

	// compute_groundtruth
	dataType = getEnvOrDefault("DATA_TYPE", "float")
	distFn = getEnvOrDefault("DIST_FN", "l2")
	baseFile = getEnvOrDefault("BASE_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")
	queryFile = getEnvOrDefault("QUERY_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fbin")
	gtFile = getEnvOrDefault("GT_FILE", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_query_learn_gt100")
	K = getEnvOrDefault("K", "100")

	// BuildDiskIndex
	indexPathPrefix = getEnvOrDefault("INDEX_PATH_PREFIX", "/home/zjlab/zyg/DiskANN/build/data/sift/disk_index_sift_learn_R32_L50_A1.2")

	// SearchDiskIndex
	resultK = getEnvOrDefault("RESULT_K", "10")
	L = getEnvOrDefault("L", "10")
	resultPath = getEnvOrDefault("RESULT_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/res")
	numNodesToCache = getEnvOrDefault("NUM_NODES_TO_CACHE", "10000")

	// 初始化原始数据和构建索引
	LearnVecToBin()
	LearnBiludIndex()
	healthState = 1
}

func main() {
	flag.Parse()

	router := gin.Default()
	//router.POST("/VecToBin", postVecToBin)
	router.POST("/SearchDiskIndex", postSearchDiskIndex)
	router.GET("/Healthy", GetHealthState)
	router.Run(":18180")
}

// 初始化接口

func LearnVecToBin() {
	err := FvecToBin(binPath, LfvecPath, LfbinPath)
	if err != nil {
		panic("LearnVec to Bin failed!")
	}
}

func LearnBiludIndex() {
	err := BuildDiskIndex(binPath, dataType, distFn, LfbinPath, indexPathPrefix)
	if err != nil {
		panic("build index failed!")
	}
}

func GetHealthState(c *gin.Context) {
	if healthState == 1 {
		c.IndentedJSON(http.StatusOK, "索引构建完毕")
	} else {
		c.IndentedJSON(http.StatusProcessing, "索引构建中")
	}
}

// 查询接口

func postVecToBin(c *gin.Context) {

	var vec2bin VecToBin

	// Call BindJSON to bind the received JSON
	if err := c.BindJSON(&vec2bin); err != nil {
		return
	}

	err := FvecToBin(binPath, vec2bin.Fvec, vec2bin.Fbin)
	if err != nil {
		return
	}
	c.IndentedJSON(http.StatusOK, vec2bin.Fbin)
}

func postSearchDiskIndex(c *gin.Context) {
	start := time.Now()
	// 2.vec to bin
	var vec2bin VecToBin

	// Call BindJSON to bind the received JSON
	if err := c.BindJSON(&vec2bin); err != nil {
		return
	}
	vec2bin.Fvec = filepath.Join(VecInitPath, vec2bin.Fvec)
	vec2bin.Fbin = strings.Replace(vec2bin.Fvec, ".vec", ".bin", -1)
	fmt.Print(vec2bin.Fvec, vec2bin.Fbin)

	err := FvecToBin(binPath, vec2bin.Fvec, vec2bin.Fbin)
	if err != nil {
		return
	}
	duration := time.Since(start)
	fmt.Println(duration)

	// 3.postComputeGroundTruth 可省略

	// 3.SearchDiskIndex

	fmt.Println("begin to search ..........")
	err, rarr := SearchDiskIndex(binPath, dataType, distFn, indexPathPrefix, vec2bin.Fbin, gtFile, resultK, L, resultPath, numNodesToCache)
	if err != nil {
		return
	}
	duration = time.Since(start)
	fmt.Println(duration)
	c.IndentedJSON(http.StatusCreated, gin.H{"data": rarr})
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
	cmd := exec.Command("sh", "-c", prg+" "+cmdString)
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
	fmt.Println("build index .....", prg+" "+cmdString)
	cmd := exec.Command("sh", "-c", prg+" "+cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Print("BuildDiskIndex:", string(stdout))

	return nil
}

// SearchDiskIndex  ./tests/search_disk_index  --data_type float --dist_fn l2 --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 --query_file data/sift/sift_query.fbin  --gt_file data/sift/sift_query_learn_gt100 -K 10 -L 10 20 30 40 50 100 --result_path data/sift/res --num_nodes_to_cache 10000
func SearchDiskIndex(bin, dataType, distFn, indexPathPrefix, queryFile, gtFile, K, L, resultPath, numNodesToCache string) (error, []string) {
	prg := bin + "search_disk_index"
	cmdString := fmt.Sprintf("--data_type " + dataType + " --dist_fn " + distFn + " --index_path_prefix " + indexPathPrefix + " --query_file " + queryFile + " -K " + K + " -L " + L + " --result_path " + resultPath + " --num_nodes_to_cache " + numNodesToCache)
	cmd := exec.Command("sh", "-c", prg+" "+cmdString)
	fmt.Println("cmd=====", prg+" "+cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return err, nil
	}

	fmt.Print("SearchDiskIndex:=======", string(stdout))

	str := `diskann answer:([\s\S]*)query result end`
	r := regexp.MustCompile(str)
	matches := r.FindStringSubmatch(string(stdout))
	fmt.Println("======", matches)
	if len(matches) < 2 {
		return err, nil
	}
	res := strings.Replace(matches[1], "\n", "", -1)
	resArr := strings.Fields(res)
	if len(resArr) == 0 {
		return err, nil
	}

	return nil, resArr
}
