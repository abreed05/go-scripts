package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"math"
	"path/filepath"
	"time"
)

type Primary struct {
	Agent Agent `yaml:"agent"`
}

// Agent
type Agent struct {
	Client Client `yaml:"client"`
	Database Database `yaml:"database"`
}

// Client
type Client struct {
	Hostname string   `yaml:"hostname"`
	DiskPath string	  `yaml:"diskpath"`
}

type Database struct {
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
	Name       string `yaml:"name"`
	Username   string `yaml:"username"`
	Password   string `yaml:"password"`
	Collection string `yaml:"collection"`
}

type Metric struct {
	Host string
	Total_Memory uint64
	Used_Memory uint64
	Free_Memory uint64
	Total_Disk_Space uint64
	Used_Disk_Space uint64
	Free_Disk_Space uint64
	CPU_Percent int
	Updated_At time.Time
}


func main() {

	// Command line option to set path to config.yml
	var configFile string
	flag.StringVar(&configFile, "config", "./config.yml", "Full path to config file")
	flag.StringVar(&configFile, "c", "./config.yml", "Full path to config file" )
	flag.Parse()

	// Read config YAML file
	filename, _ := filepath.Abs(configFile)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var config Primary
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// Mongo DB configuration
	dbHost := config.Agent.Database.Host
	dbPort := config.Agent.Database.Port
	dbName := config.Agent.Database.Name
	//dbUsername := config.Server.Database.Username
	//dbPassword := config.Server.Database.Password
	dbCollection := config.Agent.Database.Collection

	// Mongo Connection String
	mongoConnectionString := "mongodb://" + dbHost + ":" + dbPort
	clientOptions := options.Client().ApplyURI(mongoConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Metric Collection
	clientHostName := config.Agent.Client.Hostname
	dpath := config.Agent.Client.DiskPath
	collection := client.Database(dbName).Collection(dbCollection)
	updatedAt := time.Now()

	// CPU usage
	percent, err := cpu.Percent(time.Second,false)
	if err != nil {
		log.Fatal(err)
	}
	cpuPercent := int(math.Ceil(percent[0]))

	// Memory Usage
	vm, _ := mem.VirtualMemory()
	totalmem := bToMb(vm.Total)
	meminuse := bToMb(vm.Used)
	memfree := bToMb(vm.Free)

	// Disk usage
	du, _ := disk.Usage(dpath)
	duTotal := bToMb(du.Total)
	duFree := bToMb(du.Free)
	duUsed := bToMb(du.Used)

	// Send data to Mongo
	metricModel := Metric{clientHostName,totalmem, meminuse, memfree, duTotal,duUsed,duFree,cpuPercent,updatedAt}

	// Mongo query db.metrics.find({"host":"test2","totalmemory":{$exists: true},"updatedat":{$lt: ISODate("2021-06-02T17:40:00.000Z")}}, {"totalmemory": 1}).pretty()

	insertMetric, err := collection.InsertOne(context.TODO(), metricModel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertMetric.InsertedID)

}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

