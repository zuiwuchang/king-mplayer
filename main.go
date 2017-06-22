package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	var al bool
	flag.BoolVar(&al, "al", false, "add music to list")
	var ar bool
	flag.BoolVar(&ar, "ar", false, "random and loop play music by list")
	var ao bool
	flag.BoolVar(&ao, "ao", false, "order play music by list")
	var ap bool
	flag.BoolVar(&ap, "ap", false, "order and loop play music by list")

	var ext string
	flag.StringVar(&ext, "ext", "", "find music file ext")
	var src string
	flag.StringVar(&src, "src", "", "music list path or find path")
	var dist string
	flag.StringVar(&dist, "dist", "", "if set al then dist is list path")

	flag.Parse()
	if al {
		ext = strings.TrimSpace(ext)
		src = strings.TrimSpace(src)
		dist = strings.TrimSpace(dist)
		if ext == "" {
			fmt.Println("-ext not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}
		if src == "" {
			fmt.Println("-src not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}
		if dist == "" {
			fmt.Println("-dist not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}

		searchToList(ext, src, dist)
	} else if ar {
		src = strings.TrimSpace(src)
		if src == "" {
			fmt.Println("-src not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}

		playRandom(src)
	} else if ao {
		src = strings.TrimSpace(src)
		if src == "" {
			fmt.Println("-src not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}
		playOrder(src)
	} else if ap {
		src = strings.TrimSpace(src)
		if src == "" {
			fmt.Println("-src not set")
			fmt.Println("   -al -ext=ape -src=~/music -dist=~/my.list")
			return
		}

		playLoop(src)
	} else {
		flag.PrintDefaults()
	}
}
func loadList(dist string) ([]string, map[string]bool, error) {
	keys := make(map[string]bool)
	arrs := make([]string, 0, 1024)

	//讀取 已經存在的 歌曲
	f, e := os.Open(dist)
	if e == nil {
		r := bufio.NewReader(f)
		for {
			b, _, e := r.ReadLine()
			if e != nil {
				break
			}
			str := string(b)
			if _, ok := keys[str]; !ok {
				keys[str] = true
				arrs = append(arrs, str)
			}
		}
		f.Close()
	}
	return arrs, keys, e
}
func searchToList(ext, src, dist string) {
	arrs, keys, _ := loadList(dist)

	ext = strings.ToLower(ext)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	//尋找歌曲
	filepath.Walk(src, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(path), ext) {
			return nil
		}
		path, e = filepath.Abs(path)
		if e == nil {
			if _, ok := keys[path]; !ok {
				keys[path] = true
				arrs = append(arrs, path)
			}
		}
		return nil
	})

	f, e := os.Create(dist)
	if e != nil {
		fmt.Println(e)
		return
	}
	for _, str := range arrs {
		f.WriteString(str + "\n")
	}
	f.Close()
}
func playOrder(src string) {
	arrs, _, e := loadList(src)
	if e != nil {
		fmt.Println(e)
		return
	}
	n := len(arrs)
	if n == 0 {
		fmt.Println("list empty")
		return
	}

	for i := 0; i < n; i++ {
		play(arrs[i])
	}
}
func playLoop(src string) {
	arrs, _, e := loadList(src)
	if e != nil {
		fmt.Println(e)
		return
	}
	n := len(arrs)
	if n == 0 {
		fmt.Println("list empty")
		return
	}

	for {
		for i := 0; i < n; i++ {
			play(arrs[i])
		}
	}
}
func playRandom(src string) {
	arrs, _, e := loadList(src)
	if e != nil {
		fmt.Println(e)
		return
	}
	n := len(arrs)
	if n == 0 {
		fmt.Println("list empty")
		return
	}
	rand.Seed(time.Now().Unix())
	for {
		i := rand.Int() % n
		play(arrs[i])
	}
}
func play(str string) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("mplayer '%s'", str))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	e := cmd.Run()
	if e != nil {
		fmt.Println(e)
	}
}
