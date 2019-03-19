package main

import (
	"context"
	"fmt"
	"sort"
  	"bufio"
  	"os"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/google/go-github/github"
)

// LatestVersions returns a sorted slice with the highest version as its first element and the highest version of the smaller minor versions in a descending order
func LatestVersions(releases []*semver.Version, minVersion *semver.Version) []*semver.Version {
	var versionSlice []*semver.Version
	// This is just an example structure of the code, if you implement this interface, the test cases in main_test.go are very easy to run

	// use "sort" interface to sort the versions in descending order
	sort.Sort(Collection(releases))

	// remove versions smaller than minVersion
	// keep only the highest patch version for each release
	for _, version := range releases {
		if version.Compare(*minVersion) >= 0{
			if len(versionSlice) == 0 {
				versionSlice = append(versionSlice, version)
			}else {
				// Compare the minor part to the sorted "versionSlice"
				if version.Slice()[1] !=  versionSlice[len(versionSlice)-1].Slice()[1]{
					versionSlice = append(versionSlice, version)
				}
			}
		}
	}
	return versionSlice
}

// Implement sort interface
type Collection []*semver.Version
func (c Collection) Len() int {
	return len(c)
}
func (c Collection) Less(i, j int) bool {
	return c[j].LessThan(*c[i]) // sort in descending order
}
func (c Collection) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}


// Read text file line by line
func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()
  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
			lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

// split into repo, subrepo and minVer
func preprocess(lines []string) ([]string, []string, []string){
	var repo, subrepo, minVer []string
	for _, line := range lines {
		temp1 := strings.Split(line, "/")
		repo = append(repo, temp1[0])
		temp2 := strings.Split(temp1[1], ",")
		subrepo = append(subrepo, temp2[0])
		minVer = append(minVer, temp2[1])
	}
	return repo, subrepo, minVer
}

// Github Release API
func gitHubReleaseAPI(repo string, subrepo string) []*semver.Version {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.ListOptions{PerPage: 10}
	releases, _, err := client.Repositories.ListReleases(ctx, repo, subrepo, opt)

	defer func(repo string, subrepo string){
		if r := recover(); r != nil{
			fmt.Printf("The repo \"%s/%s\" is not found \n", repo, subrepo)
			fmt.Println(r)
		}
	}(repo, subrepo)

	if err != nil {
		panic(err) // is this really a good way?
	}

	allReleases := make([]*semver.Version, len(releases))
	for i, release := range releases {
		versionString := *release.TagName
		if versionString[0] == 'v' {
			versionString = versionString[1:]
		}
		allReleases[i] = semver.New(versionString)
	}
	return allReleases
}


// Here we implement the basics of communicating with github through the library as well as printing the version
// You will need to implement LatestVersions function as well as make this application support the file format outlined in the README
// Please use the format defined by the fmt.Printf line at the bottom, as we will define a passing coding challenge as one that outputs
// the correct information, including this line
func main() {
	// Read text file
	lines, err := readLines("test.txt")
	if err != nil {
		fmt.Println(err)
	}
	repo, subrepo, minVer := preprocess(lines)
	for i, _ := range lines {
		versionSlice := LatestVersions(gitHubReleaseAPI(repo[i], subrepo[i]), semver.New(minVer[i]))
		fmt.Printf("latest versions of %s/%s: %s \n", repo[i], subrepo[i], versionSlice)
	}
}
