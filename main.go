package main

import (
	"flag"
	"fmt"
	"os/exec"
)

func main() {
	binPath := flag.String("binPath", "/home/zjlab/zyg/bin/", "bin path")
	fvecPath := flag.String("fvecPath", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fvecs", "vec path")
	fbinPath := flag.String("fbinPath", "/home/zjlab/zyg/DiskANN/build/data/sift/sift_learn.fbin", "bin path")

	flag.Parse()

	fmt.Println(*binPath)
	FvecToBin(*binPath,*fvecPath,*fbinPath)

}

func FvecToBin(bin ,fvecPath,fbinPath string)  {
	prg := bin+"fvecs_to_bin"


	cmd := exec.Command(prg,fvecPath,fbinPath)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
}