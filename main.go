package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

var (
	binPath string
	fvecPath string
	fbinPath string

	dataType string
	distFn string
	baseFile string
	queryFile string
	gtFile string
	K string

	dataPath string
	indexPathPrefix string

	resultK string
	L string
	resultPath string
	numNodesToCache string
)

func init() {
	// vec_to_bin
	binPath = getEnvOrDefault("BIN_PATH", "/home/zjlab/zyg/bin/")
	fvecPath = getEnvOrDefault("FVEC_PATH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fvecs")
	fbinPath = getEnvOrDefault("FBIN_PAThH", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")

	// compute_groundtruth
	dataType = getEnvOrDefault("DATA_TYPE", "float")
	distFn = getEnvOrDefault("DIST_FN","l2")
	baseFile = getEnvOrDefault("BASE_FILE","/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")
	queryFile = getEnvOrDefault("QUERY_FILE","/home/zjlab/zyg/DiskANN/build/data/sift/sift_query.fbin")
	gtFile = getEnvOrDefault("GT_FILE","/home/zjlab/zyg/DiskANN/build/data/sift/sift_query_learn_gt100")
	K = getEnvOrDefault("K","100")

	// BuildDiskIndex
	dataPath = getEnvOrDefault("DATA_PATH","/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin")
	indexPathPrefix = getEnvOrDefault("INDEX_PATH_PREFIX","/home/zjlab/zyg/DiskANN/build/data/sift/disk_index_sift_learn_R32_L50_A1.2")

	// SearchDiskIndex
	resultK = getEnvOrDefault("RESULT_K","10")
	L = getEnvOrDefault("L","10 20 30 40 50 100")
	resultPath = getEnvOrDefault("RESULT_PATH","/home/zjlab/zyg/DiskANN/build/data/sift/res")
	numNodesToCache = getEnvOrDefault("NUM_NODES_TO_CACHE","10000")
}

func main() {
	flag.Parse()
	
	FvecToBin(binPath,fvecPath,fbinPath)
	ComputeGroundTruth(binPath,dataType,distFn,baseFile,queryFile,gtFile,K)
	BuildDiskIndex(binPath,dataType,distFn,dataPath,indexPathPrefix)
	SearchDiskIndex(binPath,dataType,distFn,indexPathPrefix,queryFile,gtFile,resultK,L,resultPath,numNodesToCache)
}

func getEnvOrDefault(env string, defaultValue string) string{
	v := os.Getenv(env)
	if v == "" {
		v = defaultValue
	}
	return v
}


func FvecToBin(bin ,fvecPath,fbinPath string)  {
	prg := bin+"fvecs_to_bin"
	cmd := exec.Command(prg,fvecPath,fbinPath)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("FvecToBin:",string(stdout))
}

//ComputeGroundTruth  ./tests/utils/compute_groundtruth  --data_type float --dist_fn l2 --base_file data/sift/sift_learn.fbin --query_file  data/sift/sift_query.fbin --gt_file data/sift/sift_query_learn_gt100 --K 100
func ComputeGroundTruth(bin ,dataType,distFn,baseFile,queryFile,gtFile ,K string )  {
	prg := bin+"compute_groundtruth"
	cmdString := fmt.Sprintf("--data_type "+dataType+" --dist_fn "+distFn+" --base_file "+baseFile+" --query_file "+queryFile+" --gt_file "+gtFile+" --K  "+K)
	fmt.Println(cmdString)
	cmd := exec.Command(prg,cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("ComputeGroundTruth:",string(stdout))
}

//BuildDiskIndex  ./tests/build_disk_index --data_type float --dist_fn l2 --data_path data/sift/sift_learn.fbin --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 -R 32 -L50 -B 0.003 -M 1
func BuildDiskIndex(bin ,dataType,distFn,dataPath ,indexPathPrefix string)  {
	prg := bin+"build_disk_index"
	subCmd := " -R 32 -L50 -B 0.003 -M 1"
	cmdString := fmt.Sprintf("--data_type "+dataType+" --dist_fn "+distFn+" --data_path "+dataPath+" --index_path_prefix "+indexPathPrefix+subCmd)
	cmd := exec.Command(prg,cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("BuildDiskIndex:",string(stdout))
}

//SearchDiskIndex  ./tests/search_disk_index  --data_type float --dist_fn l2 --index_path_prefix data/sift/disk_index_sift_learn_R32_L50_A1.2 --query_file data/sift/sift_query.fbin  --gt_file data/sift/sift_query_learn_gt100 -K 10 -L 10 20 30 40 50 100 --result_path data/sift/res --num_nodes_to_cache 10000
func SearchDiskIndex(bin,dataType,distFn,indexPathPrefix,queryFile,gtFile,K ,L,resultPath ,numNodesToCache string)  {
	prg := bin+"search_disk_index"
	cmdString := fmt.Sprintf("--data_type "+dataType+" --dist_fn "+distFn+" --index_path_prefix "+indexPathPrefix+" --query_file "+queryFile+" --gt_file "+gtFile+" -K "+K+" -L "+L+" --result_path "+resultPath+" --num_nodes_to_cache "+numNodesToCache)
	cmd := exec.Command(prg,cmdString)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print("SearchDiskIndex:",string(stdout))
}