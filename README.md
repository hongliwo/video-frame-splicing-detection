# 使用说明

核心是两

1. 使用ffmpeg去处理视频以及合成图片
2. 使用AWS GO SDK V2去下载视频以及使用Amazon Rekognition做图片推理

## 使用AWS GO SDK V2去处理S3视频文件和图像监测

[AWS Go SDK v2](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/)

[Detect Labels API](https://docs.aws.amazon.com/zh_cn/rekognition/latest/APIReference/API_DetectLabels.html)

```bash
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
#for rek
go get github.com/aws/aws-sdk-go-v2/service/rekognition
# for s3
go get github.com/aws/aws-sdk-go-v2/service/s3
```

## 具体请看main函数

> 视频需要提前上传到S3 切需要配置相关权限确保SDK能够正常执行代码

```go
func main() {

    //s3文件路径
    bucket := "s3.plaza.red"
    key := "test/test.mov"
    
    //下载文件到本地
    filename, filepath := GetVideoFromS3(bucket, key)
    
    //视频抽桢合成
    imagePath := VideoToImage(filename, filepath)
    
    // rekognition call
    DetectLabelsByRekognition(imagePath)

}

```

### 依赖安装
如果找不到ffmpeg，可以在此地址下载
https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz

### 执行结果
get object from hongliwotestbucket/video_20230310_105453.mov
CallCommand Run 参数=>  [-i /tmp/1678696848081152833.mov -t 4 -s 640x360 -r 1 /tmp/1678696848081152833/frame%d.jpg]
CallCommand Run 执行命令=>  /usr/bin/ffmpeg -i /tmp/1678696848081152833.mov -t 4 -s 640x360 -r 1 /tmp/1678696848081152833/frame%d.jpg

CallCommand Run 调用完成.....
CallCommand Run 参数=>  [-i /tmp/1678696848081152833/frame1.jpg -i /tmp/1678696848081152833/frame2.jpg -filter_complex hstack /tmp/1678696848081152833/p12.jpg]
CallCommand Run 执行命令=>  /usr/bin/ffmpeg -i /tmp/1678696848081152833/frame1.jpg -i /tmp/1678696848081152833/frame2.jpg -filter_complex hstack /tmp/1678696848081152833/p12.jpg

CallCommand Run 调用完成.....
CallCommand Run 参数=>  [-i /tmp/1678696848081152833/frame3.jpg -i /tmp/1678696848081152833/frame4.jpg -filter_complex hstack /tmp/1678696848081152833/p34.jpg]
CallCommand Run 执行命令=>  /usr/bin/ffmpeg -i /tmp/1678696848081152833/frame3.jpg -i /tmp/1678696848081152833/frame4.jpg -filter_complex hstack /tmp/1678696848081152833/p34.jpg

CallCommand Run 调用完成.....
CallCommand Run 参数=>  [-i /tmp/1678696848081152833/p12.jpg -i /tmp/1678696848081152833/p34.jpg -filter_complex vstack /tmp/1678696848081152833/p1234.jpg]
CallCommand Run 执行命令=>  /usr/bin/ffmpeg -i /tmp/1678696848081152833/p12.jpg -i /tmp/1678696848081152833/p34.jpg -filter_complex vstack /tmp/1678696848081152833/p1234.jpg

CallCommand Run 调用完成.....
read file /tmp/1678696848081152833/p1234.jpg len: 107668 
[{"Aliases":[],"Categories":[{"Name":"Home and Indoors"}],"Confidence":99.60803,"Instances":[],"Name":"Cushion","Parents":[{"Name":"Home Decor"}]},{"Aliases":[],"Categories":[{"Name":"Furniture and Furnishings"}],"Confidence":99.60803,"Instances":[],"Name":"Home Decor","Parents":[]},{"Aliases":[],"Categories":[{"Name":"Furniture and Furnishings"}],"Confidence":84.28107,"Instances":[],"Name":"Furniture","Parents":[]},{"Aliases":[],"Categories":[{"Name":"Furniture and Furnishings"}],"Confidence":84.28107,"Instances":[],"Name":"Table","Parents":[{"Name":"Furniture"}]},{"Aliases":[],"Categories":[{"Name":"Technology and Computing"}],"Confidence":79.19881,"Instances":[],"Name":"Computer Hardware","Parents":[{"Name":"Electronics"},{"Name":"Hardware"}]},{"Aliases":[],"Categories":[{"Name":"Technology and Computing"}],"Confidence":79.19881,"Instances":[],"Name":"Electronics","Parents":[]},{"Aliases":[],"Categories":[{"Name":"Furniture and Furnishings"}],"Confidence":66.76585,"Instances":[],"Name":"Couch","Parents":[{"Name":"Furniture"}]},{"Aliases":[],"Categories":[{"Name":"Person Description"}],"Confidence":65.38793,"Instances":[{"BoundingBox":{"Height":0.205918,"Left":0.49907413,"Top":0.49839428,"Width":0.10302267},"Confidence":65.38793,"DominantColors":null}],"Name":"Adult","Parents":[{"Name":"Person"}]},{"Aliases":[],"Categories":[{"Name":"Person Description"}],"Confidence":65.38793,"Instances":[{"BoundingBox":{"Height":0.205918,"Left":0.49907413,"Top":0.49839428,"Width":0.10302267},"Confidence":65.38793,"DominantColors":null}],"Name":"Man","Parents":[{"Name":"Adult"},{"Name":"Male"},{"Name":"Person"}]},{"Aliases":[{"Name":"Human"}],"Categories":[{"Name":"Person Description"}],"Confidence":65.38793,"Instances":[{"BoundingBox":{"Height":0.205918,"Left":0.49907413,"Top":0.49839428,"Width":0.10302267},"Confidence":65.38793,"DominantColors":null}],"Name":"Person","Parents":[]}]