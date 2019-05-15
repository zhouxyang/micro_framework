package route

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"route_guide/cmd"
	"route_guide/configfile"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	pb "route_guide/routeguide"
)

func init() {
	//注册服务初始化函数
	cmd.RegisterService("myService", InitServer)

}

// InitServer 初始化myService服务
func InitServer(grpcServer *grpc.Server, config *configfile.Config) error {
	srv := &GuideServer{
		RouteNotes: make(map[string][]*pb.RouteNote),
		config:     config,
	}
	err := srv.LoadFeatures(config.JSONDBFile)
	if err != nil {
		return err
	}
	pb.RegisterRouteGuideServer(grpcServer, srv)
	return nil
}

// GuideServer 自定义服务结构体
type GuideServer struct {
	savedFeatures []*pb.Feature // read-only after initialized
	config        *configfile.Config

	mu         sync.Mutex // protects RouteNotes
	RouteNotes map[string][]*pb.RouteNote
}

// printFeatures lists all the features within the given bounding Rectangle.
func printFeatures(client pb.RouteGuideClient, rect *pb.Rectangle, log *logrus.Entry) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
	}
	count := 0
	for {
		_, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", client, err)
		}
		count++
	}
	log.Infof("printFeatures finish: %d", count)
}

// GetFeature returns the feature at the given point.
func (s *GuideServer) GetFeature(ctx context.Context, point *pb.Point) (*pb.Feature, error) {
	log := cmd.GetLog(ctx)
	for _, feature := range s.savedFeatures {
		if proto.Equal(feature.Location, point) {
			log.Infof("GetFeatureSucc")
			cli, err := cmd.GetEtcdClient(s.config.EtcdHost, s.config.EtcdPort)
			if err != nil {
				log.Infof("GetEtcdClient error:%v", err)
				return feature, nil
			}
			defer cli.Close()
			conn, err := cmd.GetGrpcConn(ctx, "myService", cli, log)
			if err != nil {
				log.Infof("GetEtcdClient error:%v", err)
				return feature, nil
			}
			defer conn.Close()
			client := pb.NewRouteGuideClient(conn)
			printFeatures(client, &pb.Rectangle{
				Lo: &pb.Point{Latitude: 400000000, Longitude: -750000000},
				Hi: &pb.Point{Latitude: 420000000, Longitude: -730000000},
			}, log)
			return feature, nil
		}
	}
	// No feature was found, return an unnamed feature
	return &pb.Feature{Location: point}, nil
}

// ListFeatures lists all features contained within the given bounding Rectangle.
func (s *GuideServer) ListFeatures(rect *pb.Rectangle, stream pb.RouteGuide_ListFeaturesServer) error {
	log := cmd.GetLog(stream.Context())
	for _, feature := range s.savedFeatures {
		if inRange(feature.Location, rect) {
			if err := stream.Send(feature); err != nil {
				return err
			}
		}
	}
	log.Infof("ListFeaturesSucc")
	return nil
}

// RecordRoute records a route composited of a sequence of points.
//
// It gets a stream of points, and responds with statistics about the "trip":
// number of points,  number of known features visited, total distance traveled, and
// total time spent.
func (s *GuideServer) RecordRoute(stream pb.RouteGuide_RecordRouteServer) error {
	var pointCount, featureCount, distance int32
	var lastPoint *pb.Point
	startTime := time.Now()
	for {
		point, err := stream.Recv()
		if err == io.EOF {
			endTime := time.Now()
			return stream.SendAndClose(&pb.RouteSummary{
				PointCount:   pointCount,
				FeatureCount: featureCount,
				Distance:     distance,
				ElapsedTime:  int32(endTime.Sub(startTime).Seconds()),
			})
		}
		if err != nil {
			return err
		}
		pointCount++
		for _, feature := range s.savedFeatures {
			if proto.Equal(feature.Location, point) {
				featureCount++
			}
		}
		if lastPoint != nil {
			distance += calcDistance(lastPoint, point)
		}
		lastPoint = point
	}
}

// RouteChat receives a stream of message/location pairs, and responds with a stream of all
// previous messages at each of those locations.
func (s *GuideServer) RouteChat(stream pb.RouteGuide_RouteChatServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		key := serialize(in.Location)

		s.mu.Lock()
		s.RouteNotes[key] = append(s.RouteNotes[key], in)
		// Note: this copy prevents blocking other clients while serving this one.
		// We don't need to do a deep copy, because elements in the slice are
		// insert-only and never modified.
		rn := make([]*pb.RouteNote, len(s.RouteNotes[key]))
		copy(rn, s.RouteNotes[key])
		s.mu.Unlock()

		for _, note := range rn {
			if err := stream.Send(note); err != nil {
				return err
			}
		}
	}
}

// LoadFeatures loads features from a JSON file.
func (s *GuideServer) LoadFeatures(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "Failed to load default features")
	}
	if err := json.Unmarshal(file, &s.savedFeatures); err != nil {
		return errors.Wrap(err, "Failed to load default features")
	}
	return nil
}

func toRadians(num float64) float64 {
	return num * math.Pi / float64(180)
}

// calcDistance calculates the distance between two points using the "haversine" formula.
// The formula is based on http://mathforum.org/library/drmath/view/51879.html.
func calcDistance(p1 *pb.Point, p2 *pb.Point) int32 {
	const CordFactor float64 = 1e7
	const R float64 = float64(6371000) // earth radius in metres
	lat1 := toRadians(float64(p1.Latitude) / CordFactor)
	lat2 := toRadians(float64(p2.Latitude) / CordFactor)
	lng1 := toRadians(float64(p1.Longitude) / CordFactor)
	lng2 := toRadians(float64(p2.Longitude) / CordFactor)
	dlat := lat2 - lat1
	dlng := lng2 - lng1

	a := math.Sin(dlat/2)*math.Sin(dlat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dlng/2)*math.Sin(dlng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return int32(distance)
}

func inRange(point *pb.Point, rect *pb.Rectangle) bool {
	left := math.Min(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	right := math.Max(float64(rect.Lo.Longitude), float64(rect.Hi.Longitude))
	top := math.Max(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))
	bottom := math.Min(float64(rect.Lo.Latitude), float64(rect.Hi.Latitude))

	if float64(point.Longitude) >= left &&
		float64(point.Longitude) <= right &&
		float64(point.Latitude) >= bottom &&
		float64(point.Latitude) <= top {
		return true
	}
	return false
}

func serialize(point *pb.Point) string {
	return fmt.Sprintf("%d %d", point.Latitude, point.Longitude)
}
