package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"
)

var AWSCredentialProfile string = "gws" //这个根据你实际需要处理

func CallCommandRun(cmdName string, args []string) (string, error) {
	cmd := exec.Command(cmdName, args...)
	fmt.Println("CallCommand Run 参数=> ", args)
	fmt.Println("CallCommand Run 执行命令=> ", cmd)
	bytes, err := cmd.Output()
	if err != nil {
		fmt.Println("CallCommand Run 出错了.....", err.Error())
		fmt.Println(err)
		return "", err
	}
	resp := string(bytes)
	fmt.Println(resp)
	fmt.Println("CallCommand Run 调用完成.....")
	return resp, nil
}

/**
利用ffmpeg抽帧
*/
func VideoToImage(filename string, filepath string) string {
	//videoFile := "./test.mov"

	ffmpeg := "ffmpeg" // 命令安装的位置， 也可以直接指定绝对路径

	workdir := "/tmp/" + filename

	//创建目录
	err := os.MkdirAll(workdir, 0755)
	if err != nil {
		fmt.Println(err)
	}

	//开始ffmpeg的工作
	CallCommandRun(ffmpeg, []string{"-i", filepath, "-t", "4", "-s", "640x360", "-r", "1", workdir + "/frame%d.jpg"})
	//fmt.Println(x)
	CallCommandRun(ffmpeg, []string{"-i", workdir + "/frame1.jpg", "-i", workdir + "/frame2.jpg", "-filter_complex", "hstack", workdir + "/p12.jpg"})

	CallCommandRun(ffmpeg, []string{"-i", workdir + "/frame3.jpg", "-i", workdir + "/frame4.jpg", "-filter_complex", "hstack", workdir + "/p34.jpg"})

	CallCommandRun(ffmpeg, []string{"-i", workdir + "/p12.jpg", "-i", workdir + "/p34.jpg", "-filter_complex", "vstack", workdir + "/p1234.jpg"})

	return workdir + "/p1234.jpg"
}

// 读取文件到[]byte中
func file2Bytes(filename string) ([]byte, error) {

	// File
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// FileInfo:
	stats, err := file.Stat()
	if err != nil {
		return nil, err
	}

	// []byte
	data := make([]byte, stats.Size())
	count, err := file.Read(data)
	if err != nil {
		return nil, err
	}
	fmt.Printf("read file %s len: %d \n", filename, count)
	return data, nil
}

func GetVideoFromS3(bucket string, key string) (string, string) {

	fmt.Printf("get object from %s/%s\n", bucket, key)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.Region = "ap-southeast-1"
		//o.UseAccelerate = true
	})

	res, err := s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		panic(err)
	}

	//在临时文件目录下生成文件名称
	filename := strconv.Itoa(int(time.Now().UnixNano())) //随机数生成文件名称
	filepath := "/tmp/" + filename + ".mov"
	outFile, err := os.Create(filepath)
	// handle err
	defer outFile.Close()
	_, err = io.Copy(outFile, res.Body)
	// handle err
	return filename, filepath
}

//调用API接口

func DetectLabelsByRekognition(image string) string {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	rekClient := rekognition.NewFromConfig(cfg, func(o *rekognition.Options) {
		o.Region = "us-east-1"
	})

	//这两个参数根据实际调试参数去调整
	var ml int32 = 10
	var mc float32 = 0.75
	//两个参数调整和验证

	//read bytes
	imageContent, err := file2Bytes(image)
	if err != nil {
		panic(err)
	}

	//这里是Rekognition的调用现场
	res, err := rekClient.DetectLabels(context.TODO(), &rekognition.DetectLabelsInput{
		Image: &types.Image{
			Bytes: imageContent,
		},
		//Features:      nil,
		MaxLabels:     &ml,
		MinConfidence: &mc,
		//Settings:      nil,
	})

	if err != nil {
		panic(err)
	}

	jsonLabels, _ := json.Marshal(res.Labels)
	fmt.Printf("%s\n", jsonLabels)

	return string(jsonLabels)
}

func main() {

	//s3文件路径 := 自动推导
	bucket := "hongliwotestbucket"
	key := "video_20230310_105453.mov"

	//下载文件到本地
	filename, filepath := GetVideoFromS3(bucket, key)
	//视频抽桢合成
	imagePath := VideoToImage(filename, filepath)
	// rekognition call
	DetectLabelsByRekognition(imagePath)

}
