package bemetrics

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/oschwald/geoip2-golang"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/woozymasta/a2s/pkg/a2s"
	"github.com/woozymasta/bercon-cli/pkg/beparser"
	"github.com/woozymasta/bercon-cli/pkg/bercon"
)

// reads test data from a specified text file.
func LoadTestData(filename string) ([]byte, error) {
	data, err := os.ReadFile(filepath.Join("test_data", filename))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func initVars() (rcon_address, password, query_address string) {
	ip, ok := os.LookupEnv("DAYZ_EXPORTER_RCON_IP")
	if !ok {
		ip = "127.0.0.1"
	}
	rcon_port, ok := os.LookupEnv("DAYZ_EXPORTER_RCON_PORT")
	if !ok {
		rcon_port = "2025"
	}
	rcon_address = fmt.Sprintf("%s:%s", ip, rcon_port)

	query_port, ok := os.LookupEnv("DAYZ_EXPORTER_QUERY_PORT")
	if !ok {
		query_port = "27016"
	}
	query_address = fmt.Sprintf("%s:%s", ip, query_port)

	password, ok = os.LookupEnv("DAYZ_EXPORTER_RCON_PASSWORD")
	if !ok {
		password = ""
	}

	return
}

func getCustomLabels() Labels {
	return Labels{
		{Key: "AAA", Value: "111"},
		{Key: "BBB", Value: "222"},
	}
}

func TestPlayerMetricsFromFile(t *testing.T) {
	// load data from file
	input, err := LoadTestData("players.txt")
	if err != nil {
		t.Fatalf("Failed to load players test data: %v", err)
	}

	// parse data
	players := &beparser.Players{}
	players.Parse(input)
	if err != nil {
		t.Fatalf("Error parsing players: '%v'", err)
	}

	geoDB, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		t.Errorf("Cant open GeoDB %e", err)
	}
	defer geoDB.Close()
	players.SetCountryCode(geoDB)

	// initialize and register metrics for players
	mc := NewMetricsCollector(getCustomLabels())
	mc.InitPlayerMetrics()
	mc.RegisterMetrics()

	// update metrics with player data
	mc.UpdatePlayerMetrics(players)

	// output metrics
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}

	for _, mf := range metricsFamilies {
		fmt.Printf("Metric family: %s\n", *mf.Name)
		for _, m := range mf.Metric {
			fmt.Println(m)
		}
	}
}

func TestMetricsFromBercon(t *testing.T) {
	// initialize the parameters for connection
	address, password, query_address := initVars()

	// connect to A2S Query
	query, err := a2s.NewWithString(query_address)
	if err != nil {
		t.Fatalf("Error connecting to A2S Query: %v", err)
	}
	info, err := query.GetInfo()
	if err != nil {
		t.Fatalf("Error query A2S_INFO: %v", err)
	}

	// connect to BeRCON
	berconClient, err := bercon.Open(address, password)
	if err != nil {
		t.Fatalf("Error connecting to Bercon: %v", err)
	}
	defer berconClient.Close()

	// send the command "players" and receive a response
	data, err := berconClient.Send("players")
	if err != nil {
		t.Fatalf("Error sending 'players' command: %v", err)
	}

	// parse response
	players := beparser.NewPlayers()
	players.Parse(data)

	// initialize and register metrics for players
	mc := NewMetricsCollector(getCustomLabels())
	mc.InitServerMetrics() // A2S
	mc.InitPlayerMetrics() // RCON
	mc.RegisterMetrics()

	mc.UpdateServerMetrics(info)
	// update metrics with player data
	mc.UpdatePlayerMetrics(players)

	// output metrics
	metricsFamilies, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}

	for _, mf := range metricsFamilies {
		fmt.Printf("Metric family: %s\n", *mf.Name)
		for _, m := range mf.Metric {
			fmt.Println(m)
		}
	}
}

func TestBansMetricsFromBerconSeparateRegistry(t *testing.T) {
	// initialize the parameters for connection
	address, password, _ := initVars()

	// connect to BeRCON
	berconClient, err := bercon.Open(address, password)
	if err != nil {
		t.Fatalf("Error connecting to Bercon: %v", err)
	}
	defer berconClient.Close()

	// send the command "players" and receive a response
	data, err := berconClient.Send("bans")
	if err != nil {
		t.Fatalf("Error sending 'players' command: %v", err)
	}

	// parse response
	response := beparser.Parse(data, "bans")
	bans, ok := response.(*beparser.Bans)
	if !ok {
		t.Fatalf("Parsed data is not of type 'beparser.Bans")
	}

	geoDB, err := geoip2.Open("../beparser/GeoLite2-Country.mmdb")
	if err != nil {
		t.Errorf("Cant open GeoDB %e", err)
	}
	defer geoDB.Close()
	bans.SetCountryCode(geoDB)

	//? create new metrics registry
	reg := prometheus.NewRegistry()

	// initialize and register metrics for players
	mc := NewMetricsCollector(getCustomLabels())
	mc.InitBansMetrics()
	reg.MustRegister(mc.banGUIDTimeMetric)
	reg.MustRegister(mc.banGUIDTotal)
	reg.MustRegister(mc.banIPTimeMetric)
	reg.MustRegister(mc.banIPTotal)

	// update metrics with player data
	mc.UpdateBansMetrics(bans)

	// output metrics
	metricsFamilies, err := reg.Gather()
	if err != nil {
		t.Fatalf("Error gathering metrics: %v", err)
	}

	for _, mf := range metricsFamilies {
		fmt.Printf("Metric family: %s\n", *mf.Name)
		for _, m := range mf.Metric {
			fmt.Println(m)
		}
	}
}
