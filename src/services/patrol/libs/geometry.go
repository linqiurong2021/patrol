package libs

import (
	"errors"
	"fmt"
	"github.com/linqiurong2021/patrol/src/services/patrol/structs"
	"strconv"
	"strings"
)

const (
	Polyline = "polyline"
	Polygon = "polygon"
)


// 点
type Point struct {
	Longitude float64 `json:"longitude"`
	Latitude float64 `json:"latitude"`
	Height float32 `json:"height"`
}

type Geometry struct {

}

//
func NewGeometry() *Geometry {
	return &Geometry{}
}

func (g *Geometry) PointText2Point(pointText string) *structs.Point  {
	tmp := []byte(pointText)
	//
	data := tmp[6:len(tmp)-1]
	//
	return g.text2Point(data)[0]
}

func (g *Geometry) PointText2Points(pointText string) []*structs.Point  {
	tmp := []byte(pointText)
	//
	data := tmp[6:len(tmp)-1]
	//
	return g.text2Point(data)
}

func (g *Geometry) Polyline2Points(polyline string) []*structs.Point {
	tmp := []byte(polyline)
	//
	data := tmp[11:len(tmp)-1]
	//
	return g.text2Point(data)
}

func (g *Geometry) Polygon2Points(polygon string) []*structs.Point  {
	tmp := []byte(polygon)
	data := tmp[9:len(tmp)-1]
	//
	return g.text2Point(data);
}

// 点转文本
func (g *Geometry) Points2Text(points []*structs.Point, geometryType string) (string,error)  {
	// 判断如果点位信息不存在则直接返回点位信息不存在
	if len(points) <=0 {
		return "", errors.New("patrol polyline not found")
	}
	// 点位文本信息
	var pointsText string
	//
	for _,point := range points{
		// 判断是否是点位
		if !g.IsPoint(point) {
			return "", nil
		}
		// 连接
		pointsText += fmt.Sprintf("%f %f,", point.Longitude,point.Latitude)
	}
	// 去掉最后一个,
	pointsText = pointsText[0: len(pointsText) -1]
	//
	if geometryType == Polyline {
		//
		pointsText = fmt.Sprintf("LINESTRING(%s)", pointsText)
		//LINESTRING(118.48 24.46, 118.44 24.4, 118.4 24.48)
	}else if geometryType == Polygon {
		//
		pointsText = fmt.Sprintf("POLYGON((%s))", pointsText)
		//POLYGON((118.185824 24.489506, 118.18604 24.48974, 118.186701 24.489505, 118.186411 24.4887, 118.185585 24.489056, 118.185824 24.489506))
	}

	return pointsText, nil
}

// 判断是否是一个正常的点位
func (g *Geometry) IsPoint(point *structs.Point) bool  {
	//
	if (point.Longitude > -180 && point.Longitude < 180) && (point.Latitude >-90 && point.Latitude< 90){
		return true
	}
	return false
}
// 单个点转换
func (g *Geometry) Point2Text(point *structs.Point) (string,error) {
	//
	if !g.IsPoint(point) {
		return "", errors.New("point invalidate")
	}
	//
	pointText := fmt.Sprintf("Point(%f %f)",point.Longitude, point.Longitude)
	return pointText, nil
}




func (g *Geometry) text2Point(data []byte) []*structs.Point  {
	// data 为去除point(  ) Linestring( )  Polygon( ) 后的数据
	dataArr := strings.Split(string(data),",")
	var points []*structs.Point
	//
	for _,item :=range dataArr {
		point := strings.Fields(item)

		longitude, _:=strconv.ParseFloat(point[1],64)
		latitude, _:=strconv.ParseFloat(point[0],64)
		var height float64

		p := structs.Point{
			Latitude: longitude,
			Longitude: latitude,
			Height: height,
		}
		points = append(points,&p)
	}
	return points
}