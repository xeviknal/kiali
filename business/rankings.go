package business

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/kiali/kiali/graph"
	"github.com/kiali/kiali/graph/api"
	"github.com/kiali/kiali/graph/config/cytoscape"
	"github.com/kiali/kiali/kubernetes"
	"github.com/kiali/kiali/prometheus"
	"strconv"
)

type RankingService struct {
	k8s           kubernetes.ClientInterface
	prom          prometheus.ClientInterface
	rd 			  redis.Client
	businessLayer *Layer
}

var ctx = context.Background()

// Create Redis client
var rd = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})

type ItemWithScore struct {
	Item string
	Score float64
}

func (in *RankingService) Load(options graph.Options) (int, error) {
	// Fetch metrics from n intervals (15s, 1m, 5m)
	var cytosConfig cytoscape.Config
	code, config := api.GraphNamespaces(in.businessLayer, options)
	if code != 0 {
		return code, nil
	}

	if config != nil {
		cytosConfig = config.(cytoscape.Config)
	}

	for _, edge := range cytosConfig.Elements.Edges {
		if rt, err := strconv.ParseFloat(edge.Data.ResponseTime, 64); err == nil {
			// Slow edges
			rd.ZIncrBy(ctx, "slow-edges", rt, edge.Data.ID)
		}

		tp, tpErr := strconv.ParseFloat(edge.Data.Throughput, 64)
		rps, rpsErr := strconv.ParseFloat(edge.Data.Traffic.Rates["http"], 64)

		if tpErr == nil && rpsErr == nil {
			// Most time consuming
			rd.ZIncrBy(ctx, "mtc-edges", tp * rps, edge.Data.ID)
		}
	}

	return 0, nil
}

func (in *RankingService) getRanking(set string, max int) []ItemWithScore {
	res := rd.ZRangeWithScores(ctx, set, 0, int64(max) - 1)
	rank := res.Val()
	iws := make([]ItemWithScore, 0, len(rank))
	for _, z := range rank {
		iws = append(iws, ItemWithScore{
			Item:  z.Member.(string),
			Score: z.Score,
		})
	}
	return iws
}

func (in *RankingService) GetSlowEdges(max int) []ItemWithScore {
	return in.getRanking("slow-edges", max)
}

func (in *RankingService) GetMostTimeConsumingEdges(max int) []ItemWithScore {
	return in.getRanking("mtc-edges", max)
}
