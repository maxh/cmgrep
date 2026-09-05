import pb "github.com/maxh/cmgrep/proto"

type server struct {
      pb.UnimplementedCmgrepServer
}

func (s *server) GrepCount(ctx context.Context, req *pb.GrepCountRequest) (*pb.GrepCountResponse, error) {
      n, err := countMatches(req.GetPattern(), req.GetIgnoreCase(), req.GetExtended())
      if err != nil {
              return nil, status.Errorf(codes.InvalidArgument, "bad pattern: %v", err)
      }
      return &pb.GrepCountResponse{Count: n}
}
